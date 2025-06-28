package storage

import (
	"database/sql"
	"fmt"
	"go-http-basics/models"
	"strings"
)

func GetAllUsersFromDB(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func InsertUserToDb(db *sql.DB, user *models.User) error {
	query := "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id"
	return db.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
}

func DeleteUserFromBd(db *sql.DB, id int) (bool, error) {
	query := "DELETE FROM users WHERE id = $1"

	res, err := db.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func UpdateUserInDb(db *sql.DB, id int, input models.UpdateUserInput) (bool, error) {
	parts := []string{}
	args := []interface{}{}
	i := 1

	if input.Name != nil {
		parts = append(parts, fmt.Sprintf("name = $%d", i))
		args = append(args, *input.Name)
		i++
	}

	if input.Age != nil {
		parts = append(parts, fmt.Sprintf("age = $%d", i))
		args = append(args, *input.Age)
		i++
	}

	if len(parts) == 0 {
		return false, nil
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(parts, ", "), i)
	args = append(args, id)
	res, err := db.Exec(query, args...)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
