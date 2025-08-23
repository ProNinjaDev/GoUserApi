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

func main() {
	http.HandleFunc("/", handleRoot)

	log.Println("Запуск сервера")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
