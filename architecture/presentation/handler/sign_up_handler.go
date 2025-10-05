package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"time"
)

// Go отдаёт JSON ({status: "ok", user: {...}}), JS меняет DOM
// Потом когда эррор хэндлинг до конца будет сделан
// бэк будет передавать статус код ошибки и в нетворке
// фронт будет показывать его

func (m *MainHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Method not Allowed",
		})
		//http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var newUserRequest struct {
		Nickname  string `json:"nickname"`
		Age       string `json:"age"`
		Gender    string `json:"gender"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newUserRequest); err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   http.StatusText(http.StatusBadRequest),
		})
		return
	}
	// Parse age from date string
	age := 0
	if newUserRequest.Age != "" {
		birthDate, err := time.Parse("2006-01-02", newUserRequest.Age)
		if err == nil {
			age = int(time.Since(birthDate).Hours() / 24 / 365.25)
		}
	}

	// Basic validation
	if newUserRequest.Nickname == "" || newUserRequest.Email == "" || newUserRequest.Password == "" {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Nickname, email and password are required",
		})
		return
	}

	newUser := &models.User{
		Nickname:  newUserRequest.Nickname,
		FirstName: newUserRequest.FirstName,
		LastName:  newUserRequest.LastName,
		Age:       age,
		Gender:    newUserRequest.Gender,
		Email:     newUserRequest.Email,
		Password:  newUserRequest.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := m.service.User.Create(newUser)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}
