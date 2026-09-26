package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/entity"
)

type Repository interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUsers(ctx context.Context) ([]entity.User, error) {
	const query = `
		SELECT id, name, email, age, created_at, updated_at
		FROM users
		ORDER BY created_at ASC, id ASC
		LIMIT 100;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get users query: %w", err)
	}
	defer rows.Close()

	users := make([]entity.User, 0, 100)
	for rows.Next() {
		var current entity.User
		if err := rows.Scan(
			&current.ID,
			&current.Name,
			&current.Email,
			&current.Age,
			&current.CreatedAt,
			&current.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, current)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *repository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	const query = `
		INSERT INTO users (name, email, age)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, age, created_at, updated_at;
	`

	var created entity.User
	if err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Age).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.Age,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}
