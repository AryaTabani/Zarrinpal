package repository

import (
	"context"
	"zarrinpal/db"
	"zarrinpal/models"
)

func CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (first_name,last_name,email,password_hash) VALUES (?, ?, ?, ?)`
	_, err := db.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password_hash)
	return err
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, first_name, last_name, email, password_hash FROM users WHERE email=?`
	err := db.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password_hash,
	)
	return &user, err
}
