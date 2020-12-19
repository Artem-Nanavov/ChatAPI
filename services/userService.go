package services

import (
	"encoding/base64"
	"fmt"

	"api/database"
	"api/database/repo"
	"api/entities"

	"golang.org/x/crypto/argon2"

	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Config ...
type Config struct {
	Salt      string
	JWTSecret string
}

// UserService ...
type UserService struct {
	db     *database.Database
	config *Config
}

// NewUserService ...
func NewUserService(db *database.Database, config *Config) *UserService {
	return &UserService{
		db:     db,
		config: config,
	}
}

// Repo ...
func (u *UserService) Repo() *repo.UserRepo {
	return u.db.User()
}

func (u *UserService) hashPassword(password string) string {
	// Helper struct for password hashing
	type passwordConfig struct {
		time    uint32
		memory  uint32
		threads uint8
		keyLen  uint32
	}

	c := &passwordConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}

	hash := argon2.IDKey(
		[]byte(password),
		[]byte(u.config.Salt),
		c.time,
		c.memory,
		c.threads,
		c.keyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString([]byte(u.config.Salt))
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
}

// ComparePasswords ...
func (u *UserService) ComparePasswords(user *entities.User, pass string) bool {
	return u.hashPassword(pass) == user.Password
}

func (u *UserService) validate(user *entities.User) string {
	email := validation.Validate(user.Email, is.Email)
	username := validation.Validate(user.Username, validation.Length(3, 50))
	password := validation.Validate(user.Password, validation.Length(6, 50))

	errors := ""

	if user.Password == "" || user.Username == "" {
		errors += "Must provide all data"
	}

	if email != nil {
		errors += email.Error() + ". Got: " + user.Email + ". "
	}

	if username != nil {
		errors += username.Error() + ". Got: " + user.Username + ". "
	}

	if password != nil {
		errors += password.Error() + ". Got: " + user.Password + ". "
	}

	return errors
}

// Create ...
func (u *UserService) Create(user *entities.User) (string, error) {

	// Validation
	if message := u.validate(user); message != "" {
		return "", fmt.Errorf(message)
	}

	user.Password = u.hashPassword(user.Password)

	if err := u.db.User().Create(user); err != nil {
		return "", err
	}

	return u.GenerateToken(user)
}

// GenerateToken ...
func (u *UserService) GenerateToken(user *entities.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email

	return token.SignedString([]byte(u.config.JWTSecret))
}
