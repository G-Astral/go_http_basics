package main

import (
	"fmt"
	"net/http"

	"go-http-basics/handlers"
	"go-http-basics/routes"
)

func main() {
	err := handlers.LoadFromFile()
	if err != nil {
		fmt.Println("Ошибка загрузки users.json:", err)
	}

	http.HandleFunc("/users", routes.RouteUsers)
	http.HandleFunc("/users/", routes.RouteUsersID)

	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
