package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Id     int64
	Name   string
	Status bool
}

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

// POST /user/
func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на создание пользователя")

	var user User

	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&user)

	if err != nil {
		log.Printf("Ошибка декодирования юзера JSON: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	log.Println("Создание успешное")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GET /user/{id}
func handleUserGetById(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")

	userId, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil {
		log.Printf("Не удалось сконвертировать ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("Получен запрос на получение пользователя с ID: %d", userId)

	w.Write([]byte("GET /user/{id} принят с id = " + userIdString))

}

func main() {
	//http.HandleFunc("/", handleRoot)
	//http.HandleFunc("/user/", handleUserCreate)

	r := chi.NewRouter()
	r.Get("/", handleRoot)
	r.Post("/user", handleUserCreate)

	log.Println("Запуск сервера")

	err := http.ListenAndServe(":8081", r)

	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
