package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ProNinjaDev/GoUserApi/internal/config"
	"github.com/go-chi/chi/v5"
)

type User struct {
	Id     int64
	Name   string
	Status bool
}

type api struct {
	db *sql.DB
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
func (a *api) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на создание пользователя")

	var user User

	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&user)

	if err != nil {
		log.Printf("Ошибка декодирования юзера JSON: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	query := "INSERT INTO users (name, status) VALUES ($1, $2) RETURNING id"
	err = a.db.QueryRowContext(r.Context(), query, user.Name, user.Status).Scan(&user.Id)

	if err != nil {
		log.Printf("Не удалось вставить пользователя в БД: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Println("Создание успешное")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GET /user/{id}
func (a *api) handleUserGetById(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")

	userId, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil {
		log.Printf("Не удалось сконвертировать ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("Получен запрос на получение пользователя с ID: %d", userId)

	query := "SELECT id, name, status FROM users WHERE id = $1"

	var user User
	err = a.db.QueryRowContext(r.Context(), query, userId).Scan(&user.Id, &user.Name, &user.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Пользователь с ID = %d не найден", userId)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Printf("Не удалось найти пользователя в БД: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

// GET /user?status=true&name=alex
func (a *api) handleUserGetByFilter(w http.ResponseWriter, r *http.Request) {
	statusString := r.URL.Query().Get("status")
	name := r.URL.Query().Get("name")

	log.Printf("Получен запрос на получение списка пользователей с фильтрами: status=%s, name=%s", statusString, name)
	w.Write([]byte(" GET /user/ с фильтрами status=" + statusString + ", name=" + name))
}

// PUT /user/{id}
func (a *api) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")

	userId, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil {
		log.Printf("Не удалось сконвертировать ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на изменение пользователя с ID: %d", userId)

	w.Write([]byte("PUT /user/{id} принят с id = " + userIdString))
}

// DELETE /user/{id}
func (a *api) handleUserDelete(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")

	userId, err := strconv.ParseInt(userIdString, 10, 64)

	if err != nil {
		log.Printf("Не удалось сконвертировать ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на удаление пользователя с ID: %d", userId)

	w.Write([]byte("DELETE /user/{id} принят с id = " + userIdString))

}

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	defer db.Close()

	apiObj := &api{db: db}

	r := chi.NewRouter()
	r.Route("/user", func(r chi.Router) {
		r.Post("/", apiObj.handleUserCreate)
		r.Get("/", apiObj.handleUserGetByFilter)
		r.Get("/{id}", apiObj.handleUserGetById)
		r.Put("/{id}", apiObj.handleUserUpdate)
		r.Delete("/{id}", apiObj.handleUserDelete)
	})

	r.Get("/", handleRoot)

	log.Println("Запуск сервера")

	err = http.ListenAndServe(":8081", r)

	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
