package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (m *MainHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.GetAllCategoriesHandler(w, r)
	case http.MethodPost:
		m.CreateCategoryHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (m *MainHandler) GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset := int64(0)
	limit := int64(50) // default limit

	if offsetStr != "" {
		if parsed, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			offset = parsed
		}
	}

	if limitStr != "" {
		if parsed, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = parsed
		}
	}

	categories, err := m.service.Category.GetAll(offset, limit)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success":    true,
		"categories": categories,
	})
}

func (m *MainHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

	var categoryData struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&categoryData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if categoryData.Name == "" {
		http.Error(w, "category name is required", http.StatusBadRequest)
		return
	}

	// Note: Category creation should be implemented in the service layer
	// For now, we'll return an error indicating it's not implemented
	http.Error(w, "Category creation not implemented yet", http.StatusNotImplemented)
}
