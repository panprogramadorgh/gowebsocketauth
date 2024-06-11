package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	types "github.com/panprogramadorgh/gowebsocketauth/internal/typesutils"
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Adds the client with the new connection
	var client *types.Client = types.CreateCli(conn)
	if err := clients.AddCli(client); err != nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
			fmt.Println(err)
			return
		}
	}

	// When the handler function ends the connection it's been close
	defer func() {
		if err := clients.RmCli(client, &sessions); err != nil {
			fmt.Println(err)
			return
		}
	}()

	for {
		// Reads messages from any client
		_, msg, err := (**client).ReadMessage()
		if err != nil {
			// In case of the client closes the connection, the infinite loop will stop and the deffered function will execute removing the client session (if the cli has any session)
			fmt.Println(err)
			break
		}

		message := string(msg)

		// Client is entering a command
		if string(msg[0]) == "/" {
			var command string = strings.TrimPrefix(message, "/")
			cmdoutput, msg := HandleCommand(command, client)

			var errWhenWritingMsg error = nil
			switch cmdoutput {
			case CmdOutput.PrivateMessage:
				if err := (**client).WriteMessage(websocket.TextMessage, []byte("OK "+msg)); err != nil {
					errWhenWritingMsg = err
				}
			case CmdOutput.PublicMessage:
				for _, eachClient := range clients {
					if err := (**eachClient).WriteMessage(websocket.TextMessage, []byte("OK "+msg)); err != nil {
						errWhenWritingMsg = err
					}
				}
			case CmdOutput.Error:
				if err := (**client).WriteMessage(websocket.TextMessage, []byte("ER "+msg)); err != nil {
					errWhenWritingMsg = err
				}
			}

			if errWhenWritingMsg != nil {
				fmt.Println(errWhenWritingMsg)
				break
			}
		} else {
			// Client is entering a message

			// Cheks if the client connection has a session
			currentSession := sessions.FindSessionPerCli(client)
			if currentSession == nil {
				// Reject the command
				if err := conn.WriteMessage(websocket.TextMessage, []byte("ER [server]: you are not logged in")); err != nil {
					fmt.Println(err)
					break
				}
			} else {
				for _, eachClient := range clients {
					sessionUsrName := (*currentSession).User.Username
					msgWithSession := fmt.Sprintf("[%v]: %v", sessionUsrName, message)
					if err := (**eachClient).WriteMessage(websocket.TextMessage, []byte(msgWithSession)); err != nil {
						fmt.Println(err)
						break
					}
				}
			}

		}
	}
}
