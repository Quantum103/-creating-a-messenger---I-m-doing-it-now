package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("my-super-secret-key-change-in-production")

// Middleware для проверки JWT токена
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Токен не предоставлен", http.StatusUnauthorized)
			return
		}

		// Убираем "Bearer "
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Парсим и проверяем токен
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
			return
		}

		// Токен валидный — передаём управление следующему обработчику
		next(w, r)
	}
}

func createProxy(host string) *httputil.ReverseProxy {
	target, _ := url.Parse("http://" + host)

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// Убираем исходные заголовки, чтобы прокси установил свои
		req.Header.Del("User-Agent")
		req.Header.Del("X-Forwarded-For")
		req.Header.Del("X-Forwarded-Host")
	}

	// Обработчик ошибок
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf(" Ошибка прокси: %v", err)
		http.Error(w, "Сервис временно недоступен", http.StatusServiceUnavailable)
	}

	return proxy
}
func main() {
	r := mux.NewRouter()

	authProxy := createProxy("localhost:8081")
	userProxy := createProxy("localhost:8082")
	// регистрация
	r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
			return
		}
		authProxy.ServeHTTP(w, r)
	}).Methods("POST")

	// логин
	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
			return
		}
		authProxy.ServeHTTP(w, r)
	}).Methods("POST")

	// c middleware

	r.HandleFunc("/api/dashboard", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем данные из токена
		tokenString := r.Header.Get("Authorization")[7:]
		claims := jwt.MapClaims{}
		token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if token.Valid {
			// Добавляем данные пользователя в заголовки для микросервиса
			if userID, ok := claims["user_id"].(float64); ok {
				r.Header.Set("X-User-ID", fmt.Sprintf("%.0f", userID))
			}
			if email, ok := claims["email"].(string); ok {
				r.Header.Set("X-User-Email", email)
			}
			if username, ok := claims["username"].(string); ok {
				r.Header.Set("X-User-Email", username)
			}
		}

		userProxy.ServeHTTP(w, r)
	})).Methods("GET")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../frontend/"))))

	log.Fatal(http.ListenAndServe(":8080", r))
}
