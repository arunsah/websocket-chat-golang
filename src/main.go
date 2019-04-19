package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket" // go get github.com/gorilla/websocket
)

// TODO: The access to the clients map is not concurrency-safe.
// Try using sync.Map 
var clients = make(Map[*websocket.Conn]bool)  // global : currently connected clients
var broadcast = make(chan Message)            // global: broadcast message queue/ channel
var connectionUpgrader = websocket.Upgrader{} // normal HTTP connections to websocket

// Message will contains details of user
// Email can be used to fetch unique gravatar
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// simple static file server binded to root
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// configure websocket route
	http.HandleFunc("/ws", handleWSConnections)

	// takes messages from the broadcast channel from before and pass them to clients
	go handleMessages()
	// starts the web server on localhost and log any errors
	Host := flag.String("h", "localhost", "host IP for this web app")
	Port := flag.String("p", "8000", "web server port")
	flag.Parse()
	host := *Host
	port := ":" + *Port
	log.Println("http socket based web server is starting at ", host, port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("main:ListenAndServe", err)
	}

}

func handleWSConnections(res http.ResponseWriter, req *http.Request) {
	// will run for each web request and run indefinitely until there is some error
	// from client (we assume that the connection is closed)

	// upgrade initial GET request to a websocket
	ws, err := connectionUpgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Fatalln("handleWSConnections:upgrader,Upgrade", err)
		// close the connection before function returns
		defer ws.Close()

	}

	// registering the client in global map
	clients[ws] = true

	// infinite loop
	for {
		var msg Message
		// read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("handleWSConnections:ReadJSON:", err)
			delete(clients, ws)
			break
		}
		// send the newly received message to the broadcast message queue
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// grab the next message from the message queue
		msg := <-broadcast
		// send ot out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("handleMessages:WriteJSON", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
