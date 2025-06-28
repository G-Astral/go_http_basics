package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go-http-basics/db"
	"go-http-basics/models"
	"go-http-basics/storage"
)

var (
	users  []models.User
	nextID = 1
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	// fmt.Println("Сохраняем пользователей в файл...")
	err = SaveToFile()
	if err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	users = append(users[:index], users[index+1:]...)

	err = SaveToFile()
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Пользователь удален")
}

func PatchUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var input models.UpdateUserInput
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

	err = SaveToFile()
	if err != nil {
		http.Error(w, "Ошибка при сохранении", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Пользователь обновлён")
}

func SaveToFile() error {
	file, err := os.Create("users.json")
	// fmt.Println("Файл успешно открыт для записи")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(users)
}

func LoadFromFile() error {
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

	// fmt.Printf("Загружено пользователей: %d\n", len(users))
	// fmt.Println("Результат декодирования:", users)
	return err
}

func UsersDbHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	users, err := storage.GetAllUsersFromDB(db.DB)
	if err != nil {
		http.Error(w, "Ошибка чтения из БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&users)
}

func UserDbCreateHandler(w http.ResponseWriter, r *http.Request) {
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

	if user.Name == "" {
		http.Error(w, "Имя не может быть пустым", http.StatusBadRequest)
		return
	}

	if user.Age <= 0 {
		http.Error(w, "Возраст должен быть положительным числом", http.StatusBadRequest)
		return
	}

	err = storage.InsertUserToDb(db.DB, &user)
	if err != nil {
		log.Println("Ошибка при вставке в БД:", err)
		http.Error(w, "Ошибка записи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UserDbDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/usersDB/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	flag, err := storage.DeleteUserFromBd(db.DB, id)
	if err != nil {
		http.Error(w, "Ошибка удаления из базы данных", http.StatusInternalServerError)
		return
	}

	if !flag {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Пользователь удален")
	}
}

func UserDbPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/usersDB/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var input models.UpdateUserInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	flag, err := storage.UpdateUserInDb(db.DB, id, input)
	if err != nil {
		http.Error(w, "Ошибка в обновлении базы данных", http.StatusInternalServerError)
		return
	}

	if !flag {
		http.Error(w, "Нет данных для обновления", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Пользователь изменен")
}
