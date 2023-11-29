// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	clientsMu sync.RWMutex
	clients   = make(map[string]*websocket.Conn)
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	fmt.Println("Connecting: ", username)
	clientsMu.RLock()
	if _, ok := clients[username]; ok {
		log.Println("User already in: ", username)
		return
	}
	clientsMu.RUnlock()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		conn.Close()
		clientsMu.Lock()
		delete(clients, username)
		clientsMu.Unlock()
	}()

	clientsMu.Lock()
	clients[username] = conn
	clientsMu.Unlock()
	broadcastMessageToUsers("System", []byte(fmt.Sprintf("%s joined.", username)))
	sendAllCurrentUser(username)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			broadcastMessageToUsers("System", []byte(fmt.Sprintf("%s dropped.", username)))
			return
		}
		message := string(p)

		fmt.Printf("User %s: %s\n", username, message)

		sendMessageToUser(username, messageType, p)
	}
}

func broadcastMessageToUsers(sender string, message []byte) {
	wg := sync.WaitGroup{}
	for client := range clients {
		if client == sender {
			continue
		}
		c := client
		wg.Add(1)
		go func() {
			sendMessageToUser(sender, 1, []byte(fmt.Sprintf("%s:%s", c, message)))
			wg.Done()
		}()
	}
	wg.Wait()
}

func sendAllCurrentUser(username string) {
	var allUsers []string
	clientsMu.RLock()
	for client := range clients {
		allUsers = append(allUsers, client)
	}
	clientsMu.RUnlock()
	sendMessageToUser("System", 1, []byte(fmt.Sprintf("%s:Online[%s]", username, strings.Join(allUsers, ","))))
}

func sendMessageToUser(sender string, messageType int, message []byte) {
	fmt.Println(string(message))
	parts := strings.Split(string(message), ":")
	if len(parts) < 2 {
		broadcastMessageToUsers(sender, message)
		return
	}
	fmt.Println(parts[0], parts[1])
	receiver := parts[0]
	content := parts[1]
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	if targetConn, ok := clients[receiver]; ok {
		if err := targetConn.WriteMessage(messageType, []byte(fmt.Sprintf("From %s: %s", sender, content))); err != nil {
			log.Println(err)
		}
	} else {
		notFoundMessage := fmt.Sprintf("User %s not found or not connected", receiver)
		if err := clients[sender].WriteMessage(messageType, []byte(notFoundMessage)); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
