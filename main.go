package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type User struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
	Age		int		`json:"age"`
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

func saveToFile() error {
	file, err := os.Create("users.json")
	fmt.Println("Файл успешно открыт для записи")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(users)
}

func main()  {
	http.HandleFunc("/user", createUserHandler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
