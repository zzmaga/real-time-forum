package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"time"
)

type PostRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category []string `json:"category"`
}

func (m *MainHandler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.DisplayPostsHandler(w, r)
	case http.MethodPost:
		m.CreatePostHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (m *MainHandler) DisplayPostsHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	_, err := m.service.Session.GetByUuid(authHeader)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	posts, err := m.service.Post.GetAll(0, 0)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Мне нужно имя автора и потом категории
	// Могу сам добавить или ты. Думаю ты можешь создать стракт в GetAll
	/*
			type Post struct {
		    	models.Post
		    	Author string `json:"author"`
				Categories []string `json:"category"`
			}
			и просто будешь возвращать его.
			Либо сделаешь тут как тебе виднее
	*/
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (m *MainHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	session, err := m.service.Session.GetByUuid(authHeader)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var postData PostRequest
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	newPost := &models.Post{
		Title:     postData.Title,
		Content:   postData.Content,
		UserId:    session.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// У нас оказывается нет таблицы с категориями в базе
	// добавишь таблицу с категориями связи между пост айди и категориями
	_, err = m.service.Post.Create(newPost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Post created successfully"})
}
