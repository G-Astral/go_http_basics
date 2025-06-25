package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"go-http-basics/models"
)



var (
	users 	[]models.User
	nextID 	int = 1
)

func createUserHandler(w http.ResponseWriter, r * http.Request)  {
	// Если заданный метод не является методом создания юзера - выходим
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
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
	fmt.Printf("Загружено пользователей: %d\n", len(users))
	fmt.Println("Результат декодирования:", users)
	return err
}

func main()  {
	err := loadFromFile()
	if err != nil {
		fmt.Println("Ошибка загрузки users.json:", err)
	}

	http.HandleFunc("/user", createUserHandler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
