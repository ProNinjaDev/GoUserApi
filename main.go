package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ProNinjaDev/GoUserApi/internal/config"
	"github.com/ProNinjaDev/GoUserApi/internal/user/handler"
	"github.com/ProNinjaDev/GoUserApi/internal/user/repository"
	"github.com/ProNinjaDev/GoUserApi/internal/user/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос")

	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Printf("Ошибка кодирования JSON")
	}
}

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userHandler.RegisterRoutes(r)
	r.Get("/", handleRoot)

	log.Println("Запуск сервера")

	err = http.ListenAndServe(":8081", r)

	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
