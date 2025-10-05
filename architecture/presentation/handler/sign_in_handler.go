package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/architecture/service/user"
)

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
