package main

import (
	"log"
	"net/http"
	"real-time-forum/architecture/repository"
	"real-time-forum/architecture/service"
	"real-time-forum/architecture/web/handler"
	"real-time-forum/database"

	_ "github.com/mattn/go-sqlite3"
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

	// Start WebSocket message broadcaster
	go handler.HandleMessages()

	log.Println("Server starting on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
