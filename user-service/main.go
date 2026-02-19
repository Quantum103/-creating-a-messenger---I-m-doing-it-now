package main

import (
	"log"
	"net/http"
	"user-service/database"
	"user-service/handlers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConfig := database.DBConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "",
		DBName:   "micro",
	}

	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	r := http.NewServeMux()
	// Get запрос - userHandler.go
	r.HandleFunc("/dashboard", handlers.DashboardHandler(db))

	// POST запрос - postHandler.go
	r.HandleFunc("/api/posts", handlers.PostHandler(db))

	log.Println("User Service запущен на порту 8082")
	http.ListenAndServe(":8082", r)
}
