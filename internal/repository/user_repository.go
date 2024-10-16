package repository

import (
	"ItDevTest/internal/helpers"
	"ItDevTest/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) returning id`

	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		fmt.Errorf("error getting user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) Authenticate(email, password string) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.IsBlocked {
		return nil, errors.New("user is blocked due to too many failed login attempts")
	}

	if user == nil || !helpers.CheckPasswordHash(password, user.Password) {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= 3 {
			user.IsBlocked = true
			r.UpdateUser(user) // Обновите информацию о пользователе в БД
		}
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (r *UserRepository) UpdateUserTokens(userId int, verificationToken, resetPasswordToken string, resetPasswordExpires time.Time) error {
	query := `UPDATE users SET verification_token = $1, reset_password_token = $2, reset_password_expires = $3 WHERE id = $4`
	_, err := r.db.Exec(query, verificationToken, resetPasswordToken, resetPasswordExpires, userId)
	if err != nil {
		return fmt.Errorf("error updating user tokens: %v", err)
	}
	return nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET name = $1, password = $2, failed_login_attempts = $3, is_blocked = $4 WHERE id = $5`
	_, err := r.db.Exec(query, user.Name, user.Password, user.FailedLoginAttempts, user.IsBlocked, user.Id)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}
