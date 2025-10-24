package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"zarrinpal/models"
	"zarrinpal/repository"

	"golang.org/x/crypto/bcrypt"
)

//var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

var (
	ErrUserExists = errors.New("a user with this email already exists")
)

func RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error) {
	_, err := repository.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		return nil, ErrUserExists
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("database error checking user: %w", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	newUser := &models.User{
		FirstName:     payload.FirstName,
		LastName:      payload.LastName,
		Email:         payload.Email,
		Password_hash: string(hashedPassword),
	}

	err = repository.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	return newUser, nil
}
