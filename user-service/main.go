package main

import (
	"log"
	"net/http"

	"user-service/database"
	"user-service/handlers"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	r.HandleFunc("/api/register", handlers.HandleRegister(db)).Methods("POST")

	log.Println("✅ User Service запущен на порту 8081")
	log.Fatal(http.ListenAndServe(":8081", r))

}
