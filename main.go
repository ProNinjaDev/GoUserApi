package main

import (
	"encoding/json"
	"log"
	"net/http"
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

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/user/", handleUserCreate)

	log.Println("Запуск сервера")

	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
