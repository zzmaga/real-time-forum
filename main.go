package main

import (
	"log"
	"net/http"
	"real-time-forum/architecture/repository"
	"real-time-forum/architecture/service"
	"real-time-forum/architecture/web/handler"
	"real-time-forum/database"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	broadcast = make(chan WebSocketMessage)
)

func main() {
	db, err := database.InitDB("./forum.db")
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()
	repo := repository.NewRepo(db)
	srvc := service.NewService(repo)
	mainHandler := handler.NewMainHandler(srvc)
	router := mainHandler.InitRoutes()
	log.Println("Server starting on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//func usersHandler(w http.ResponseWriter, r *http.Request) { ... }
//func privateMessagesHandler(w http.ResponseWriter, r *http.Request) { ... }
//func wsHandler(w http.ResponseWriter, r *http.Request) { ... }
