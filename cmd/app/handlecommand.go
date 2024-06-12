package main

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	types "github.com/panprogramadorgh/gowebsocketauth/internal/typesutils"
)

var CmdOutput types.CmdOutputStatus = types.CmdOutputStatus{
	PrivateMessage: 0, PublicMessage: 1, RemoveClient: 2, Error: 3,
}

func HandleCommand(command string, client *types.Client) (int, string) {
	if command == "exit" {
		/* Command used for closing the session and the connection */
		outputmsg := "[server]: you have closed the connection"
		if err := clients.RmCli(client, &sessions, &outputmsg); err != nil {
			outputmsg := "[server]: " + err.Error()
			return CmdOutput.Error, outputmsg
		}
		return CmdOutput.RemoveClient, outputmsg
	} else if command == "list" {
		/* Command used for showing all the client connections */
		outputmsg := clients.GetClients()
		return CmdOutput.PrivateMessage, outputmsg

	} else if strings.Split(command, " ")[0] == "login" {
		/* Command used for loggin */

		// Checks if the client is logged in
		sessionIsAtive := sessions.SessionExistsPerCli(client)
		if sessionIsAtive {
			outputmsg := "[server]: you are already logged in"
			return CmdOutput.Error, outputmsg
		}
		// There aren't any session active for the client conn
		args := strings.Split(command, " ")
		if len(args) != 3 {
			outputmsg := "[server]: invalid command"
			return CmdOutput.Error, outputmsg
		}
		usrname := args[1]
		password := args[2]

		sessionIsActive := sessions.SessionExistsPerUsrname(usrname)
		if sessionIsActive {
			outputmsg := "[server]: session already active in other client"
			return CmdOutput.Error, outputmsg
		}
		// Checks if the user exists
		authenticatedUser := users.AuthUsr(usrname, password)

		// Adds a new session if auth
		if authenticatedUser != nil {
			newSession := &types.Session{Client: client, User: authenticatedUser}
			if err := sessions.AddSession(newSession, clients); err != nil {
				return CmdOutput.Error, "[server]: " + err.Error()
			}
			// Warns all the clients that a new client has joined
			outputmsg := "[server]: " + usrname + " joined to the server"
			return CmdOutput.PublicMessage, outputmsg
		}
		outputmsg := "[server]: invalid credentials"
		return CmdOutput.Error, outputmsg
	} else if command == "logout" {
		currentSession := sessions.FindSessionPerCli(client)
		if currentSession == nil {
			// There aren't any session active so the client can't logout
			outputmsg := "[server]: you are not logged in"
			return CmdOutput.Error, outputmsg
		}
		// Removes the session
		outputmsg := "[server]: " + currentSession.User.Username + " left the server"
		if err := sessions.RmSession(currentSession, clients); err != nil {
			return CmdOutput.Error, "[server]: " + err.Error()
		}
		return CmdOutput.PublicMessage, outputmsg
	} else if command == "sessions" {
		if len(sessions) < 1 {
			outputmsg := "[server]: there aren't any sessions active"
			return CmdOutput.PrivateMessage, outputmsg
		}
		// Lists all the sessions
		outputmsg := "[server]: sessions:\n"
		for _, eachSession := range sessions {
			outputmsg += fmt.Sprintf("%s - %v\n", eachSession.User.Username, (**eachSession.Client).RemoteAddr().String())
		}
		return CmdOutput.PrivateMessage, outputmsg
	} else if strings.Split(command, " ")[0] == "register" {
		currentSession := sessions.SessionExistsPerCli(client)
		if currentSession {
			outputmsg := "[server]: session already active"
			return CmdOutput.Error, outputmsg
		}
		// The client hasn't any session active
		args := strings.Split(command, " ")
		if len(args) != 3 {
			outputmsg := "[server]: invalid command"
			return CmdOutput.Error, outputmsg
		}
		// Creates the new user and appends it into the clients slice
		usrname := args[1]
		passwd := args[2]

		if users.UsrExistsPerUsrname(usrname) {
			outputmsg := "[server]: user already exists"
			return CmdOutput.Error, outputmsg
		}

		newUser := &types.User{
			Username: usrname,
			Password: passwd,
		}
		if err := users.AddUsr(newUser); err != nil {
			return CmdOutput.Error, "[server]: " + err.Error()
		}

		// Informs the client about the new user

		outputmsg := "[server]: " + usrname + " joined to the server"

		// Logs the usr in automatically
		newSession := &types.Session{
			User:   newUser,
			Client: client,
		}
		if err := sessions.AddSession(newSession, clients); err != nil {
			return CmdOutput.Error, "[server]: " + err.Error()
		}

		return CmdOutput.PublicMessage, outputmsg

	} else if strings.Split(command, " ")[0] == "murder" {
		args := strings.Split(command, " ")
		if len(args) != 3 {
			outputmsg := "[server]: invalid command"
			return CmdOutput.Error, outputmsg
		}
		usrname := args[1]
		passwd := args[2]

		if sessions.SessionExistsPerUsrname(usrname) {
			outputmsg := "[server]: cannot murder user with session active"
			return CmdOutput.Error, outputmsg
		}

		authenticatedUser := users.AuthUsr(usrname, passwd)
		if authenticatedUser != nil {
			if err := users.RmUsr(authenticatedUser); err != nil {
				return CmdOutput.Error, "[server]: " + err.Error()
			}
			outputmsg := "[server]: user murdered successfully"
			return CmdOutput.PrivateMessage, outputmsg
		}
		outputmsg := "[server]: invalid credentials"
		return CmdOutput.Error, outputmsg
	} else if command == "whoami" {
		currentSession := sessions.FindSessionPerCli(client)
		if currentSession == nil {
			outputmsg := "[server]: you are not logged in"
			return CmdOutput.Error, outputmsg
		}
		outputmsg := "[server]: you are " + currentSession.User.Username
		return CmdOutput.PrivateMessage, outputmsg
	} else if strings.Split(command, " ")[0] == "tell" {
		currentSession := sessions.FindSessionPerCli(client)
		if currentSession == nil {
			outputmsg := "[server]: you are not logged in"
			return CmdOutput.Error, outputmsg
		}

		args := strings.Split(command, " ")
		if len(args) < 3 {
			outputmsg := "[server]: invalid command"
			return CmdOutput.Error, outputmsg
		}

		usrReceiverName := args[1]
		msg := fmt.Sprintf("[%s tell you]: %s", currentSession.User.Username, strings.Join(args[2:], " "))

		usr := sessions.FindSessionPerUsrname(usrReceiverName)
		if usr == nil {
			outputmsg := "[server]: " + usrReceiverName + " doesn't have a session"
			return CmdOutput.Error, outputmsg
		}

		if err := (**usr.Client).WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			return CmdOutput.Error, "[server]: " + err.Error()
		}

		outputmsg := "[server]: message has been sent to " + usrReceiverName
		return CmdOutput.PrivateMessage, outputmsg
	} else {
		// The command doesn't exist
		outputmsg := "[server]: unknown command"
		return CmdOutput.Error, outputmsg
	}
}
