package repository

import (
	"context"
	"zarrinpal/db"
	"zarrinpal/models"
)

func CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `INSERT INTO payments (amount, description, authority, status) VALUES (?, ?, ?, ?)`
	_, err := db.DB.ExecContext(ctx, query, payment.Amount, payment.Description, payment.Authority, "PENDING")
	return err
}

func GetPaymentByAuthority(ctx context.Context, authority string) (*models.Payment, error) {
	var payment models.Payment
	query := `SELECT id, amount, description, status FROM payments WHERE authority = ?`
	err := db.DB.QueryRowContext(ctx, query, authority).Scan(&payment.ID, &payment.Amount, &payment.Description, &payment.Status)
	if err != nil {
		return nil, err
	}
	payment.Authority = authority
	return &payment, nil
}

func UpdatePayment(ctx context.Context, authority string, refID string, status string) error {
	query := `UPDATE payments SET ref_id = ?, status = ? WHERE authority = ?`
	_, err := db.DB.ExecContext(ctx, query, refID, status, authority)
	return err
}
