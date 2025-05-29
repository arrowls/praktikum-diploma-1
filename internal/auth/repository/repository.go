package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/auth/entity"
	autherrors "github.com/arrowls/praktikum-diploma-1/internal/auth/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *pgx.Conn
}

func (r *AuthRepository) AddUser(ctx context.Context, username, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	_, err = r.db.Exec(ctx, `
		insert into users (username, password) values ($1, $2)
	`, username, hashedPassword)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Printf(pgErr.ConstraintName)
		if pgErr.ConstraintName == "users_username_key" {
			return nil, fmt.Errorf("%w %s", autherrors.ErrUsernameExists, err.Error())
		}

	}

	if err != nil {
		return nil, fmt.Errorf("%v %w", err, apperrors.ErrBadRequest)
	}
	return r.GetUser(ctx, username, password)
}

func (r *AuthRepository) GetUser(ctx context.Context, username, password string) (*entity.User, error) {
	var id uuid.UUID
	var userPassword string
	err := r.db.
		QueryRow(ctx, `
			select id, password from users where username = $1
	    `, username).
		Scan(&id, &userPassword)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%v %w", err, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%v %w", err, apperrors.ErrUnknown)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password provided %w", apperrors.ErrUnauthorized)
	}

	return &entity.User{
		Username: username,
		ID:       id,
	}, nil
}
