package main

import (
	"fmt"
	"net/http"
	"go-http-basics/handlers"
)

func routeUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlers.CreateUserHandler(w, r)
	case http.MethodGet:
		handlers.GetUsersHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func routeUsersID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.GetUserByIdHandler(w, r)
	case http.MethodDelete:
		handlers.DeleteUserHandler(w, r)
	case http.MethodPatch:
		handlers.PatchUserHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	err := handlers.LoadFromFile()
	if err != nil {
		fmt.Println("Ошибка загрузки users.json:", err)
	}

	http.HandleFunc("/users", routeUsers)

	http.HandleFunc("/users/", routeUsersID)

	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
