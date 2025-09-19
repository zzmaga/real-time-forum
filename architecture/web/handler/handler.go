package handler

import (
	"net/http"
	"real-time-forum/architecture/service"
)

type MainHandler struct {
	service *service.Service
}

const (
	IndexFile      = "./templates/index.html"
	StaticFilesDir = "./static"
)

func NewMainHandler(service *service.Service) *MainHandler {
	return &MainHandler{service: service}
}

func (m *MainHandler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	// 1. index.html (одна страница)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, IndexFile)
	})

	// 2. static (css, js)
	fsStatic := http.FileServer(http.Dir(StaticFilesDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	// 3. API (JSON only)
	mux.HandleFunc("/api/signup", m.SignUpHandler)
	mux.HandleFunc("/api/signin", m.SignInHandler)
	mux.HandleFunc("/api/signout", m.SignOutHandler)

	mux.HandleFunc("/api/posts", m.PostsHandler)
	mux.HandleFunc("/api/posts/", m.ViewPostHandler)
	mux.HandleFunc("/api/posts/vote/", m.PostVoteHandler)
	mux.HandleFunc("/api/comments", m.CommentsHandler)
	mux.HandleFunc("/api/comments/vote/", m.CommentVoteHandler)
	mux.HandleFunc("/api/categories", m.CategoriesHandler)

	mux.HandleFunc("/api/messages", m.MessagesHandler) // WebSocket или JSON
	mux.HandleFunc("/api/users", m.UsersHandler)
	mux.HandleFunc("/api/users/profile", m.UserProfileHandler)
	mux.HandleFunc("/ws", m.WebSocketHandler)

	return mux
}
