package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"karyawan-app/internal/middleware"
	"karyawan-app/internal/models"
	"karyawan-app/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest adalah struktur untuk request login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse adalah struktur untuk response login
type LoginResponse struct {
	Token string             `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Login menangani request login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetProfile menangani request untuk mendapatkan profil user yang sedang login
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Dapatkan user ID dari context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	user, err := h.authService.GetUserProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterRoutes mendaftarkan route untuk autentikasi
func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	// Public routes
	router.HandleFunc("/auth/login", h.Login).Methods("POST")

	// Protected routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)
	authRouter.HandleFunc("/profile", h.GetProfile).Methods("GET")
}
