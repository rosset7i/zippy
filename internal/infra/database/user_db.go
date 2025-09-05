package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/zippy/internal/entity"
)

type User struct {
	DB *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *User {
	return &User{
		DB: db,
	}
}

func (u *User) Create(user *entity.User) error {
	_, err := u.DB.Exec(context.Background(),
		"INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)",
		user.Id,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (u *User) FetchByEmail(email string) (*entity.User, error) {
	var user entity.User
	rows := u.DB.QueryRow(context.Background(), "SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1", email)
	err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	return &user, err
}
