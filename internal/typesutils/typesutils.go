package types

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

/* ------------------------------------------------------ */

type CmdOutputStatus struct {
	PrivateMessage int
	PublicMessage  int
	RemoveClient   int
	Error          int
}

/* ------------------------------------------------------ */

type Client *websocket.Conn

type Session struct {
	User   *User
	Client *Client
}

type User struct {
	Username string
	Password string
}

/* ------------------------------------------------------ */

type Clients []*Client

type Sessions []*Session

type Users []*User

type WebsocketObject interface {
	Client | Session | User
}

/* ------------------------------------------------------ */

func IndexOfWO[T WebsocketObject](object *T, objects []*T) int {
	for i, eachObj := range objects {
		if object == eachObj {
			return i
		}
	}
	return -1
}

func WOExists[T WebsocketObject](object *T, objects []*T) bool {
	return IndexOfWO[T](object, objects) > -1
}

func (usrs Users) IndexOfUsrPerUsrname(usrname string) int {
	usr := usrs.FindUsrPerUsrname(usrname)
	if usr == nil {
		return -1
	} else {
		return IndexOfWO[User](usr, usrs)
	}
}

func (usrs Users) FindUsrPerUsrname(usrname string) *User {
	for _, eachUsr := range usrs {
		if (*eachUsr).Username == usrname {
			return eachUsr
		}
	}
	return nil
}

func (sessions Sessions) FindSessionPerCli(client *Client) *Session {
	for _, eachSession := range sessions {
		if eachSession.Client == client {
			return eachSession
		}
	}
	return nil
}

func (sessions Sessions) FindSessionPerUsrname(usrname string) *Session {
	for _, eachSession := range sessions {
		if eachSession.User.Username == usrname {
			return eachSession
		}
	}
	return nil
}

func (usrs Users) UsrExistsPerUsrname(usrname string) bool {
	return usrs.IndexOfUsrPerUsrname(usrname) > -1
}

func (sessions Sessions) SessionExistsPerCli(client *Client) bool {
	return sessions.FindSessionPerCli(client) != nil
}

func (sessions Sessions) SessionExistsPerUsrname(usrname string) bool {
	return sessions.FindSessionPerUsrname(usrname) != nil
}

func CreateCli(conn *websocket.Conn) *Client {
	client := Client(conn)
	return &client
}

func (clients Clients) GetClients() string {
	if len(clients) < 1 {
		return "there aren't any clients"
	}
	outputmsg := "clients:\n"
	for _, eachClient := range clients {
		outputmsg += fmt.Sprintf("%v\n", (**eachClient).RemoteAddr().String())
	}
	return outputmsg
}

func (clients *Clients) AddCli(cli *Client) error {
	cliExists := WOExists[Client](cli, *clients)
	if cliExists {
		return errors.New("client connection already exist")
	} else {
		*clients = append(*clients, cli)
		return nil
	}
}

func (clients *Clients) RmCli(cli *Client, sessions *Sessions, byeMsg *string) error {
	cliExists := WOExists[Client](cli, *clients)
	if !cliExists {
		return errors.New("client connection doesn't exist")
	}
	// Removes the client from slice
	for i, eachCli := range *clients {
		if eachCli == cli {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
			break
		}
	}
	// Removes any session asociated to the client
	for i, eachSession := range *sessions {
		if eachSession.Client == cli {
			*sessions = append((*sessions)[:i], (*sessions)[i+1:]...)
			break
		}
	}

	// Writing byeMsg
	if byeMsg != nil {
		if err := (**cli).WriteMessage(websocket.TextMessage, []byte(*byeMsg)); err != nil {
			return err
		}
	}

	return (**cli).Close() // <- Closes the connection
}

func (sessions Sessions) GetSessions() string {
	if len(sessions) < 1 {
		return "there aren't any sessions active"
	}
	listOfSessions := "sessions:\n"
	for _, eachSession := range sessions {
		listOfSessions += eachSession.User.Username + " - " + (**eachSession.Client).RemoteAddr().String() + "\n"
	}
	return listOfSessions
}

func (sessions *Sessions) AddSession(s *Session, clients Clients) error {
	exists := WOExists[Session](s, *sessions)

	sCliIsUsed := sessions.SessionExistsPerCli(s.Client)

	usrnameExists := sessions.SessionExistsPerUsrname(s.User.Username)

	if exists {
		return errors.New("session already exists")
	} else if sCliIsUsed {
		return errors.New("cannot add new session with existing connection")
	} else if usrnameExists {
		return errors.New("cannot create new session with non unique name")
	} else {
		*sessions = append(*sessions, s)
		return nil
	}
}

func (sessions *Sessions) RmSession(s *Session, clients Clients) error {
	sExists := WOExists[Session](s, *sessions)
	if !sExists {
		return errors.New("session does't exist")
	} else {
		for i, eachSession := range *sessions {
			if eachSession == s {
				*sessions = append((*sessions)[:i], (*sessions)[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (usrs *Users) AddUsr(usr *User) error {
	exists := WOExists[User](usr, *usrs)
	usrnameExists := usrs.UsrExistsPerUsrname(usr.Username)
	if exists {
		return errors.New("user already exists")
	} else if usrnameExists {
		return errors.New("usrname for new user is not unique")
	} else {
		*usrs = append(*usrs, usr)
		return nil
	}
}

func (usrs *Users) RmUsr(usr *User) error {
	exists := WOExists[User](usr, *usrs)
	if !exists {
		return errors.New("user doesn't exists")
	} else {
		usrIndex := IndexOfWO[User](usr, *usrs)
		*usrs = append((*usrs)[:usrIndex], (*usrs)[usrIndex+1:]...)
		return nil
	}
}

func (usrs Users) AuthUsr(usrname string, passwd string) *User {
	usr := usrs.FindUsrPerUsrname(usrname)
	if usr != nil && usr.Password == passwd {
		return usr
	} else {
		return nil
	}
}

/* ------------------------------------------------------ */
