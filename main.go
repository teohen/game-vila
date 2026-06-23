package main

import (
	"encoding/json"
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/game"
	"github/teohen/mgm-tto/state"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	g         game.Game
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Init() {
	g = game.New()
}

func Quit() {
	fmt.Println("QUINT")
}

func Tick() {
	ticker := time.NewTicker(time.Duration(cnts.TickInterval) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		sim := g.ServerTick()
		dto := state.FromSimulation(sim)
		data, err := json.Marshal(dto)
		if err != nil {
			log.Println("JSON error:", err)
			continue
		}

		clientsMu.Lock()
		for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("Write error:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		clientsMu.Unlock()
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	log.Println("Client connected!")

	go func() {
		defer func() {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			g.Commands = append(g.Commands, string(message))
		}
	}()
}

func main() {
	Init()
	go Tick()
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
