package handler

import (
	"encoding/json"
	"net/http"
)

func (m *MainHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.GetAllUsersHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
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

	// For now, we'll return users with messages
	// In a real implementation, you might want to return all users
	// and handle online status separately
	users, err := m.service.PrivateMessage.GetUsersWithMessages(1) // This should be the current user ID
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (m *MainHandler) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
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
	case http.MethodGet:
		m.GetUserProfileHandler(w, r, session.UserID)
	case http.MethodPut:
		m.UpdateUserProfileHandler(w, r, session.UserID)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) GetUserProfileHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	user, err := m.service.User.GetByID(userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Don't return password in response
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"user":    user,
	})
}

func (m *MainHandler) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	var updateData struct {
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
		Gender    string `json:"gender"`
		Password  string `json:"password"` // Optional - only if user wants to change it
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get current user
	user, err := m.service.User.GetByID(userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Update fields if provided
	if updateData.Nickname != "" {
		user.Nickname = updateData.Nickname
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.Age > 0 {
		user.Age = updateData.Age
	}
	if updateData.Gender != "" {
		user.Gender = updateData.Gender
	}
	if updateData.Password != "" {
		user.Password = updateData.Password
	}

	err = m.service.User.Update(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Don't return password in response
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Profile updated successfully",
		"user":    user,
	})
}
