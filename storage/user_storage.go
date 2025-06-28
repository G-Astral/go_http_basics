package storage

import (
	"database/sql"
	"go-http-basics/models"
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
