package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   "localhost:8081",
		})
		proxy.ServeHTTP(w, r)
	}).Methods("POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../frontend/"))))

	log.Println("🚀 API Gateway запущен на порту 8080")
	log.Println("   Главная:      http://localhost:8080/")
	log.Println("   Регистрация:  http://localhost:8080/register.html")
	log.Fatal(http.ListenAndServe(":8080", r))
}
