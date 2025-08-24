package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ProNinjaDev/GoUserApi/internal/user"
	"github.com/ProNinjaDev/GoUserApi/internal/user/service"
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
