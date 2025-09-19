package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"time"
)

func (m *MainHandler) CommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.CreateCommentHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	var commentData struct {
		PostID  int64  `json:"post_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&commentData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newComment := &models.PostComment{
		PostId:    commentData.PostID,
		UserId:    session.UserID,
		Content:   commentData.Content,
		CreatedAt: time.Now(),
	}

	_, err = m.service.PostComment.Create(newComment)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Comment created successfully"})
}
