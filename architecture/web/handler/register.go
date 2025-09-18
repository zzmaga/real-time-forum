package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/models"
	"time"
)

// Go отдаёт JSON ({status: "ok", user: {...}}), JS меняет DOM

// SignInHandler - POST /api/signin
func (m *MainHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Login    string `json:"loginId"` // nickname или email
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	usr, err := m.service.User.GetByNicknameOrEmail(creds.Login)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	ok, err := usr.CompareHashAndPassword(creds.Password)
	if err != nil || !ok {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	session, err := m.service.Session.Record(usr.ID)
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
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
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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
	newUser := &models.User{
		Nickname:  newUserRequest.Nickname,
		FirstName: newUserRequest.FirstName,
		LastName:  newUserRequest.LastName,
		Age:       10,
		Gender:    newUserRequest.Gender,
		Email:     newUserRequest.Email,
		Password:  newUserRequest.Password,
		CreatedAt: time.Now(),
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

}
