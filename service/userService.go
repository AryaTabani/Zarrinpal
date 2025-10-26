package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"zarrinpal/models"
	"zarrinpal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

var (
	ErrUserExists         = errors.New("a user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
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
func LoginUser(ctx context.Context, payload *models.LoginPayload) (string, error) {
	user, err := repository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(payload.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	return generateUserToken(user.ID)
}
func generateUserToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}

func GetPaymentsHistory(ctx context.Context, userID int) ([]models.Payment, error) {
	return repository.GetPaymentsHistory(ctx, userID)
}

func UpdateProfile(ctx context.Context, userID int64, payload *models.UpdateProfilePayload) error {
	return repository.UpdateUser(ctx, userID, payload)
}
