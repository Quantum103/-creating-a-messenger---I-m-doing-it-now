package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("GET /dashboard", dashboardHandler)

	log.Println("👤 User Service запущен на порту 8082")
	http.ListenAndServe(":8082", r)
}

/*

какой JSON формат ожидается от сервера в функции dashboard

{
    "user": {
        "id": 1,
        "name": "Имя пользователя",
        "status": "Статус",
        "avatar": "A",
        "avatarColor": "#8A5E3C",
        "avatarColor2": "#6D819A",
        "stats": {
            "friends": 156,
            "posts": 47,
            "followers": 234
        },
        "info": {
            "location": "Местоположение",
            "birthday": "Дата рождения",
            "work": "Работа"
        }
    }
}

{
    "posts": [
        {
            "id": 1,
            "username": "Имя пользователя",
            "content": "Текст поста",
            "time": "2 часа назад",
            "likes": 24,
            "comments": 5,
            "shares": 2,
            "liked": false
        }
    ]
}

// Запрос:
{
    "content": "Текст нового поста"
}

// Ответ:
{
    "id": 10,
    "username": "Имя пользователя",
    "content": "Текст нового поста",
    "time": "Только что",
    "likes": 0,
    "comments": 0,
    "shares": 0
}

*/
