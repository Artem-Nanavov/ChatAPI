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
		&user.IsOnline,
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
		&user.IsOnline,
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

// GetAll ...
func (u *UserRepo) GetAll() ([]*entities.User, error) {
	rows, err := u.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}

	var users []*entities.User
	for rows.Next() {
		var (
			id       int
			email    string
			username string
			password string
			isOnline bool
		)

		if err := rows.Scan(&id, &email, &username, &password, &isOnline); err != nil {
			return nil, err
		}
		user := &entities.User{
			ID:       id,
			Email:    email,
			Username: username,
			IsOnline: isOnline,
		}

		users = append(users, user)
	}
	return users, nil
}

// GoOnline ...
func (u *UserRepo) GoOnline(user *entities.User) error {
	_, err := u.db.Exec(`UPDATE users SET is_online=true WHERE id=$1`, user.ID)
	if err != nil {
		return err
	}
	user.IsOnline = true
	return nil
}

// GoOfline ...
func (u *UserRepo) GoOfline(user *entities.User) error {
	_, err := u.db.Exec(`UPDATE users SET is_online=false WHERE id=$1`, user.ID)
	if err != nil {
		return err
	}
	user.IsOnline = false
	return nil
}
