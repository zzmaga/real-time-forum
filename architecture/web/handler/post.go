package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"strconv"
	"strings"
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
	// Create post
	postID, err := m.service.Post.Create(newPost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Add categories to post
	if len(postData.Category) > 0 {
		err = m.service.Category.AddToPostByNames(postData.Category, postID)
		if err != nil {
			// Log error but don't fail the post creation
			// In production, you might want to handle this differently
			http.Error(w, "Post created but failed to add categories", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Post created successfully"})
}

func (m *MainHandler) ViewPostHandler(w http.ResponseWriter, r *http.Request) {
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

	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/posts/"))
	if err != nil {
		http.Error(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		m.GetPostHandler(w, r, int64(postId))
	case http.MethodPut:
		m.UpdatePostHandler(w, r, int64(postId), session.UserID)
	case http.MethodDelete:
		m.DeletePostHandler(w, r, int64(postId), session.UserID)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) GetPostHandler(w http.ResponseWriter, r *http.Request, postId int64) {
	post, err := m.service.Post.GetByID(postId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	comments, err := m.service.PostComment.GetAllByPostID(postId, 0, 0)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "post": post, "comments": comments})
}

func (m *MainHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request, postId int64, userID int64) {
	// Get the post first to check ownership
	post, err := m.service.Post.GetByID(postId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user owns the post
	if post.UserId != userID {
		http.Error(w, "You can only edit your own posts", http.StatusForbidden)
		return
	}

	var updateData struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Category []string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Update post fields
	if updateData.Title != "" {
		post.Title = updateData.Title
	}
	if updateData.Content != "" {
		post.Content = updateData.Content
	}
	post.UpdatedAt = time.Now()

	err = m.service.Post.Update(post)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Post updated successfully"})
}

func (m *MainHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request, postId int64, userID int64) {
	// Get the post first to check ownership
	post, err := m.service.Post.GetByID(postId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user owns the post
	if post.UserId != userID {
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	err = m.service.Post.DeleteByID(postId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Post deleted successfully"})
}

func (m *MainHandler) PostVoteHandler(w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodPost:
		m.CreatePostVoteHandler(w, r, session.UserID)
	case http.MethodDelete:
		m.DeletePostVoteHandler(w, r, session.UserID)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) CreatePostVoteHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	var voteData struct {
		PostID int64 `json:"post_id"`
		Vote   int8  `json:"vote"` // 1 for like, -1 for dislike
	}

	if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if voteData.Vote != 1 && voteData.Vote != -1 {
		http.Error(w, "vote must be 1 or -1", http.StatusBadRequest)
		return
	}

	vote := &models.PostVote{
		PostId:    voteData.PostID,
		UserId:    userID,
		Vote:      voteData.Vote,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := m.service.PostVote.Record(vote)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Vote recorded"})
}

func (m *MainHandler) DeletePostVoteHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	postIDStr := strings.TrimPrefix(r.URL.Path, "/api/posts/vote/")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	// Get user's vote for this post
	vote, err := m.service.PostVote.GetPostUserVote(userID, postID)
	if err != nil {
		http.Error(w, "vote not found", http.StatusNotFound)
		return
	}

	err = m.service.PostVote.DeleteByID(vote.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Vote deleted"})
}
