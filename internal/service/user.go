package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rtsoy/todo-app/internal/model"
	"github.com/rtsoy/todo-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	tokenTTL   = 24 * time.Hour
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserServicer {
	return &UserService{
		repository: repository,
	}
}

func (u UserService) CreateUser(user model.CreateUserDTO) (uuid.UUID, error) {
	// Email validation
	if ok := isEmailValid(user.Email); !ok {
		return uuid.Nil, errors.New("email is not valid")
	}

	// Username validation
	if ok := isUsernameValid(user.Username); !ok {
		return uuid.Nil, errors.New("username is not valid")
	}

	// Password Validation
	if ok := isPasswordValid(user.Password); !ok {
		return uuid.Nil, errors.New("password is not valid")
	}

	// Password hashing
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return uuid.Nil, errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	id, err := u.repository.Create(user)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return uuid.Nil, errors.New("email is already taken")
		}

		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return uuid.Nil, errors.New("username is already taken")
		}

		return uuid.Nil, err
	}

	return id, nil
}

func (u UserService) GenerateToken(email, password string) (string, error) {
	user, err := u.repository.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("wrong credentials")
		}

		return "", err
	}

	if ok := checkPasswordHash(password, user.PasswordHash); !ok {
		return "", errors.New("wrong credentials")
	}

	claims := jwt.MapClaims{
		"userID":  user.ID,
		"expires": time.Now().UTC().Add(tokenTTL),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func (u UserService) ParseToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Enforces that the password must be at least 8 characters long
func isPasswordValid(password string) bool {
	return len(password) > 8
}

// Ensures that the input string contains only letters (uppercase and lowercase),
// digits, underscores, and hyphens, and it must be at least 3 characters long.
func isUsernameValid(username string) bool {
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_-]{3,}$")
	return usernameRegex.MatchString(username)
}

// Ensures that the email address contains at least 3 characters before the '@' symbol
// and follows the standard email format with a domain containing at least 2 characters after the dot.
func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9._%+-]{3,}@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	return emailRegex.MatchString(email)
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}
