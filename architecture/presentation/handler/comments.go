package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"strconv"
	"strings"
	"time"
)

func (m *MainHandler) CommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.CreateCommentHandler(w, r)
	case http.MethodDelete:
		m.DeleteCommentHandler(w, r)
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

func (m *MainHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get comment ID from query parameter
	commentIDStr := r.URL.Query().Get("id")
	if commentIDStr == "" {
		http.Error(w, "comment ID is required", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	// Get the comment first to check ownership
	comment, err := m.service.PostComment.GetByID(commentID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user owns the comment
	if comment.UserId != session.UserID {
		http.Error(w, "You can only delete your own comments", http.StatusForbidden)
		return
	}

	err = m.service.PostComment.DeleteByID(commentID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Comment deleted successfully"})
}

func (m *MainHandler) CommentVoteHandler(w http.ResponseWriter, r *http.Request) {
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
		m.CreateCommentVoteHandler(w, r, session.UserID)
	case http.MethodDelete:
		m.DeleteCommentVoteHandler(w, r, session.UserID)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) CreateCommentVoteHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	var voteData struct {
		CommentID int64 `json:"comment_id"`
		Vote      int8  `json:"vote"` // 1 for like, -1 for dislike
	}

	if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if voteData.Vote != 1 && voteData.Vote != -1 {
		http.Error(w, "vote must be 1 or -1", http.StatusBadRequest)
		return
	}

	vote := &models.PostCommentVote{
		CommentId: voteData.CommentID,
		UserId:    userID,
		Vote:      voteData.Vote,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := m.service.PostCommentVote.Record(vote)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Vote recorded"})
}

func (m *MainHandler) DeleteCommentVoteHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	commentIDStr := strings.TrimPrefix(r.URL.Path, "/api/comments/vote/")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	// Get user's vote for this comment
	vote, err := m.service.PostCommentVote.GetCommentUserVote(userID, commentID)
	if err != nil {
		http.Error(w, "vote not found", http.StatusNotFound)
		return
	}

	err = m.service.PostCommentVote.DeleteByID(vote.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Vote deleted"})
}
