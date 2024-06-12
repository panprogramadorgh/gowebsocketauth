package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/panprogramadorgh/gowebsocketauth/internal/fileutils"
	types "github.com/panprogramadorgh/gowebsocketauth/internal/typesutils"

	"github.com/gorilla/websocket"
)

/* ------------------------------------------------------ */

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/* ------------------------------------------------------ */

/*
`clients` Stores each client connection
NOTE: It's important to remove the client connection once it's closed. It is also important to remove the session.
*/

var clients types.Clients
var sessions types.Sessions
var users types.Users = types.Users{
	{
		Username: "server",
		Password: "revres",
	},
}

func main() {
	var staticsPath string
	if len(os.Args) < 2 {
		staticsPath = "./internal/fileutils/views"
		fmt.Println("must include dir path for static client files")
		fmt.Printf("server will take default static files for client (%s)\n", staticsPath)
	} else {
		staticsPath = os.Args[1]
	}

	http.HandleFunc("/echo", WsHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "document")
		document, err := fileutils.ReadFile(staticsPath + "/" + "index.html")
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, document)
	})

	http.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "script")
		javascript, err := fileutils.ReadFile(staticsPath + "/" + "main.js")
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, javascript)
	})

	var port int = 3000
	fmt.Println("Server running on", port)
	errRaisingServer := http.ListenAndServe(fmt.Sprint("0.0.0.0:", port), nil)
	if errRaisingServer != nil {
		panic(errRaisingServer)
	}
}
