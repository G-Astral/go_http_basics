package routes

import (
	"go-http-basics/handlers"
	"net/http"
)

func RouteUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlers.CreateUserHandler(w, r)
	case http.MethodGet:
		handlers.GetUsersHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func RouteUsersID(w http.ResponseWriter, r *http.Request) {
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

func RouteUsersDB(w http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case http.MethodGet:
		handlers.UsersDbHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
