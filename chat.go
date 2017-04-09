package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
)

const (
	listenAddr = "localhost:9000" // server address
)

var (
	pwd, _        = os.Getwd()
	fs            = http.FileServer(http.Dir("/Users/rasyidmujahid/Repo/go/src/github.com/rasyidmujahid/chat5/public"))
	JSON          = websocket.JSON                // codec for JSON
	Message       = websocket.Message             // codec for string, []byte
	ActiveClients = make(map[*websocket.Conn]int) // map containing clients
)

// Initialize handlers and websocket handlers
func init() {
	http.Handle("/", fs)
	http.Handle("/ws", websocket.Handler(SockServer))
}

// WebSocket server to handle chat between clients
func SockServer(ws *websocket.Conn) {
	var err error
	var clientMessage string

	// cleanup on server side
	defer func() {
		if err = ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	client := ws.Request().RemoteAddr
	log.Println("Client connected:", client)

	ActiveClients[ws] = 0
	log.Println("Number of clients connected ...", len(ActiveClients))

	// for loop so the websocket stays open otherwise
	// it'll close after one Receieve and Send
	for {
		if err = Message.Receive(ws, &clientMessage); err != nil {
			// If we cannot Read then the connection is closed
			log.Println("Websocket Disconnected waiting", err.Error())
			// remove the ws client conn from our active clients
			delete(ActiveClients, ws)
			log.Println("Number of clients still connected ...", len(ActiveClients))
			return
		}

		for cs, _ := range ActiveClients {
			if err = Message.Send(cs, clientMessage); err != nil {
				// Send it out to every client that is currently connected
				log.Println("Could not send message to ", cs, err.Error())
			}
		}
	}
}

func main() {
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
