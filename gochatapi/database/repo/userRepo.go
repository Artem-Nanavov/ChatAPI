package repo

import (
	"api/entities"
	"database/sql"
)

// UserRepo ...
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo ...
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindByID ...
func (u *UserRepo) FindByID(id int) (*entities.User, error) {
	var user entities.User

	if err := u.db.QueryRow(`SELECT * FROM users WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail ...
func (u *UserRepo) FindByEmail(email string) (*entities.User, error) {
	var user entities.User

	if err := u.db.QueryRow(`SELECT * FROM users WHERE email = $1`, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

// Create ...
func (u *UserRepo) Create(user *entities.User) error {
	if _, err := u.FindByEmail(user.Email); err != nil && err != sql.ErrNoRows {
		return err
	}
	return u.db.QueryRow(`INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id`,
		user.Email,
		user.Username,
		user.Password,
	).Scan(&user.ID)
}
