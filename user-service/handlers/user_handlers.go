package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// структура ДЛЯ сервера
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// структура ИЗ сервера
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: метод=%s, URL=%s", r.Method, r.URL.Path)

		if r.Method != http.MethodPost {
			http.Error(w, "Вход не разрешён", http.StatusMethodNotAllowed)
			return
		}

		// ограничение по памяти 1МБ
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Password == "" || req.Username == "" {
			http.Error(w, "Все поля обязательны", http.StatusBadRequest)
			return
		}
		if len(req.Password) < 4 {
			http.Error(w, "пароль слишком маленький", http.StatusBadRequest)
			return
		}
		hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "ошибка обработки пароля", http.StatusBadRequest)
			return
		}

		// готовим вставку в SQL
		query := `
    INSERT INTO users (username, email, password, created_at, updated_at) 
    VALUES (?, ?, ?, NOW(), NOW())
`

		result, err := db.Exec(query, req.Username, req.Email, string(hashPass))
		if err != nil {
			log.Printf("Ошибка при сохранении: %v", err)

			if strings.Contains(err.Error(), "Duplicate entry") {
				http.Error(w, "Пользователь с таким email или логином уже существует", http.StatusConflict)
				return
			}
			http.Error(w, "Ошибка сохранения пользователя", http.StatusInternalServerError)
			return
		}
		// получ последнего ID
		userID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Ошибка получения ID ", http.StatusInternalServerError)
			return
		}

		response := UserResponse{
			ID:        userID,
			Username:  req.Username,
			Email:     req.Email,
			CreatedAt: time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "пользователь создан",
			"user":    response,
		})
	}
}
