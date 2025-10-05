package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"real-time-forum/architecture/service/user"
	"time"
)

// Go отдаёт JSON ({status: "ok", user: {...}}), JS меняет DOM
// Потом когда эррор хэндлинг до конца будет сделан
// бэк будет передавать статус код ошибки и в нетворке
// фронт будет показывать его
// SignInHandler - POST /api/signin
func (m *MainHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Method not Allowed",
		})
		//http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Login    string `json:"loginId"` // nickname или email
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"token":   "",
			"error":   "invalid request body",
		})
		return
	}
	usr, err := m.service.User.GetByNicknameOrEmail(creds.Login)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"token":   "",
			"error":   "user not found",
		})
		return
		//http.Error(w, "user not found", http.StatusUnauthorized)
		//return
	}
	ok, err := user.CompareHashAndPassword(creds.Password)
	if err != nil || !ok {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"token":   "",
			"error":   "invalid credentials",
		})
		return
		//http.Error(w, "invalid password", http.StatusUnauthorized)
		//return
	}
	session, err := m.service.Session.Record(usr.ID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"token":   "",
			"error":   "failed to create a record",
		})
		return
	}
	// куку ставим, чтобы JS мог узнать, что пользователь залогинен
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.Uuid,
		Path:     "/",
		Expires:  session.ExpiredAt,
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"token":   session.Uuid,
		/*"user": map[string]any{
			"id":       usr.ID,
			"nickname": usr.Nickname,
		},*/
	})
}

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

/*
Потом закинешь сайнаут в другой файл куда надо.
Я не понял почему ты Signin сюда закинул если он делает заход,
а если SignUp и является регистрацией
*/
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
