package data

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
	"tranquara.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Password struct {
	plainText      *string
	hashedPassword []byte
}
type User struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	FullName  string    `json:"full_name"`
	Age       int8      `json:"age"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	query := `INSERT INTO users
				(email, password_hash, full_name, age, activated)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING user_id, email, full_name, age, activated, created_at
				`

	arg := []any{user.Email, user.Password.hashedPassword, user.FullName, user.Age, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, arg...).Scan(
		&user.UserID,
		&user.Email,
		&user.FullName,
		&user.Age,
		&user.Activated,
		&user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `
			SELECT user_id, email, full_name, age, activated  from users
			WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var user User
	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.FullName,
		&user.Age,
		&user.Activated)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
        UPDATE users 
        SET email = $2, full_name = $3, password_hash = $4, age = $5, activated = $6
        WHERE user_id = $1
        RETURNING user_id, email, full_name, age, activated `
	args := []interface{}{
		user.UserID,
		user.Email,
		user.FullName,
		user.Email,
		user.Password.hashedPassword,
		user.Age,
		user.Activated,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID,
		&user.Email,
		&user.FullName,
		&user.Age,
		&user.Activated)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (p *Password) Set(plainTextPassword string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)

	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hashedPassword = hash

	return nil
}

func (p Password) Match(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hashedPassword, []byte(plainTextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, bcrypt.ErrMismatchedHashAndPassword
		default:
			return false, err
		}

	}
	return true, nil

}

func ValidateEmail(v *validator.Validator, email string) {
	_, err := mail.ParseAddress(email)
	v.Check(email != "", "email", "must be provided")
	v.Check(err == nil, "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FullName != "", "name", "must be provided")
	v.Check(len(user.FullName) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plainText != nil {
		ValidatePasswordPlaintext(v, *user.Password.plainText)
	}

	if user.Password.hashedPassword == nil {
		panic("missing password hash for user")
	}
}
