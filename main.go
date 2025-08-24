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

	query := "SELECT id, name, status FROM users WHERE 1=1"

	args := []any{}
	argCnt := 1

	if name != "" {
		query += " AND name LIKE $" + strconv.Itoa(argCnt)
		args = append(args, name)
		argCnt++
	}

	if statusString != "" {
		status, err := strconv.ParseBool(statusString)
		if err != nil {
			log.Printf("Неверное значение для status: %s", statusString)
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		query += " AND status = $" + strconv.Itoa(argCnt)
		args = append(args, status)
		argCnt++
	}

	rows, err := a.db.QueryContext(r.Context(), query, args...)

	if err != nil {
		log.Printf("Не удалось выполнить запрос к БД: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Name, &user.Status); err != nil {
			log.Printf("Не удалось сканировать из БД: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Произошла ошибка во время цикла по строкам из БД: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

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

	var user User

	jsonDecoder := json.NewDecoder(r.Body)
	err = jsonDecoder.Decode(&user)

	if err != nil {
		log.Printf("Ошибка декодирования юзера JSON: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	query := "UPDATE users SET name = $1, status = $2 WHERE id = $3"
	result, err := a.db.ExecContext(r.Context(), query, user.Name, user.Status, userId)

	if err != nil {
		log.Printf("Не удалось обновить пользователя в БД: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Не удалось получить количество обновленных строк: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("Не удалось найти пользователя с id = %d", userId)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Id = userId
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	log.Printf("Успешно обновился пользователь с id = %d", userId)
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

	query := "DELETE FROM users WHERE id = $1"
	result, err := a.db.ExecContext(r.Context(), query, userId)
	if err != nil {
		log.Printf("Не удалось удалить пользователя из БД: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Не удалось получить количество обновленных строк: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("Не удалось найти пользователя с id = %d", userId)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Пользователь с id %d успешно удален", userId)
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
