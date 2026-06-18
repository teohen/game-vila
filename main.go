package main

import (
	"fmt"
	"github/teohen/mgm-tto/game"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	g game.Game
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
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		state := g.Update()
		fmt.Println(state)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected!")
	done := make(chan struct{})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			close(done)
			break
		}

		g.Commands = append(g.Commands, string(message))
	}
}

func main() {
	Init()
	go Tick()
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
