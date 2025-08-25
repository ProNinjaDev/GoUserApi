package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ProNinjaDev/GoUserApi/internal/user"
	"github.com/ProNinjaDev/GoUserApi/internal/user/service"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service service.Service
}

func NewUserHandler(s service.Service) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на создание пользователя")

	var u user.User

	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&u)

	if err != nil {
		log.Printf("Ошибка декодирования юзера JSON: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	err = h.service.Create(r.Context(), &u)
	if err != nil {
		log.Printf("Не удалось вставить пользователя в БД: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Println("Создание успешное")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")

	userId, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil {
		log.Printf("Не удалось сконвертировать ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("Получен запрос на получение пользователя с ID: %d", userId)

	u, err := h.service.GetByID(r.Context(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Не удалось получить id пользователя: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}
