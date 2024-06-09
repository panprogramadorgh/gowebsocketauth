package types

import "github.com/gorilla/websocket"

type User struct {
	Username string
	Password string
}

type Session struct {
	User User
	Conn *websocket.Conn
}
