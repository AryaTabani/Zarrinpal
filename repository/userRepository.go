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

func GetPaymentsHistory(ctx context.Context, userId int) ([]models.Payment, error) {
	query := `SELECT id, amount, description, status FROM payments WHERE user_id = ?`
	rows, err := db.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []models.Payment

	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.Amount, &p.Description, &p.Status); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
