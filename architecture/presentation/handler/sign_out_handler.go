package handler

import (
	"encoding/json"
	"net/http"
)

func (m *MainHandler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Method not Allowed",
		})
		return
	}

	// Get session from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "No active session",
		})
		return
	}

	// Get session by UUID first
	session, err := m.service.Session.GetByUuid(cookie.Value)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Session not found",
		})
		return
	}

	// Delete session from database
	err = m.service.Session.Delete(session.ID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Failed to delete session",
		})
		return
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Successfully signed out",
	})
}
