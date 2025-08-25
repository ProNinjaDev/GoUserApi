package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ProNinjaDev/GoUserApi/internal/config"
	"github.com/ProNinjaDev/GoUserApi/internal/user/handler"
	"github.com/ProNinjaDev/GoUserApi/internal/user/repository"
	"github.com/ProNinjaDev/GoUserApi/internal/user/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
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

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var cfg config.Config

	err = viper.Unmarshal(&cfg)

	db, err := config.ConnectDatabase(cfg)
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

	err = http.ListenAndServe(cfg.ServerPort, r)

	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
