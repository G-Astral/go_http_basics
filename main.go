package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
	Age		int		`json:"age"`
}

type UpdateUserInput struct {
	Name	*string	`json:"name"`
	Age		*int	`json:"age"`
}

var (
	users 	[]User
	nextID 	int = 1
)

func createUserHandler(w http.ResponseWriter, r * http.Request)  {
	// Если заданный метод не является методом создания юзера - выходим
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка чтения JSON", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Имя не может быть пустым", http.StatusBadRequest)
		return
	}

	if user.Age <= 0 {
		http.Error(w, "Возраст должен быть положительным числом", http.StatusBadRequest)
		return
	}

	user.ID = nextID
	nextID++
	users = append(users, user)

	fmt.Println("Сохраняем пользователей в файл...")
	err = saveToFile()
	if err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func saveToFile() error {
	file, err := os.Create("users.json")
	fmt.Println("Файл успешно открыт для записи")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(users)
}

func loadFromFile() error {
	file, err := os.Open("users.json")
	
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	
	err = json.NewDecoder(file).Decode(&users)

	maxID := 0
	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	nextID = maxID + 1

	fmt.Printf("Загружено пользователей: %d\n", len(users))
	fmt.Println("Результат декодирования:", users)
	return err
}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "Пользователь не найден", http.StatusNotFound)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	index := -1
	for i, user := range users {
		if user.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	users = append(users[:index], users[index + 1:]...)

	err = saveToFile()
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Пользователь удален")
}

func patchUserHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodPatch {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var input UpdateUserInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	if input.Name == nil && input.Age == nil {
		http.Error(w, "Нет данных для обновления", http.StatusBadRequest)
		return
	}

	found := false
	for i, user := range users {
		if user.ID == id {
			if input.Name != nil {
				users[i].Name = *input.Name
			}
			if input.Age != nil {
				users[i].Age = *input.Age
			}
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	err = saveToFile()
	if err != nil {
		http.Error(w, "Ошибка при сохранении", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Пользователь обновлён")
}

func main()  {
	err := loadFromFile()
	if err != nil {
		fmt.Println("Ошибка загрузки users.json:", err)
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createUserHandler(w, r)
		case http.MethodGet:
			getUsersHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUserByIdHandler(w, r)
		case http.MethodDelete:
			deleteUserHandler(w, r)
		case http.MethodPatch:
			patchUserHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
