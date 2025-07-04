package sqlite

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"sso/internal/domain/models"
	"sso/internal/repository"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}

func (s *Storage) Save(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	stmt, err := s.DB.Prepare("insert into users (email, pass_hash) values (?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, email, passwordHash)
	if err != nil {
		if ext, ok := err.(interface{ ExtendedCode() int }); ok && ext.ExtendedCode() == 2067 {
			return 0, repository.ErrUserExists
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	stmt, err := s.DB.Prepare("select id, email, pass_hash from users where email = ?")
	if err != nil {
		return models.User{}, err
	}
	row := stmt.QueryRowContext(ctx, email)
	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil

}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	var IsAdmin bool
	stmt, err := s.DB.Prepare("select is_admin from users where id = ?")
	if err != nil {
		return false, err
	}
	err = stmt.QueryRowContext(ctx, id).Scan(&IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.ErrUserNotFound
		}
		return false, err
	}
	return IsAdmin, nil

}

func (s *Storage) App(ctx context.Context, AppId int) (models.App, error) {
	var app models.App
	stmt, err := s.DB.Prepare("select id, name, secret from apps where id = ?")
	if err != nil {
		return models.App{}, err
	}
	err = stmt.QueryRowContext(ctx, AppId).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, repository.ErrAppNotFound
		}
		return models.App{}, err
	}
	return app, nil

}
