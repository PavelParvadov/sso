package auth

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"sso/internal/domain/models"
	"sso/internal/repository"
	"sso/pkg/jwt"
	"time"
)

var (
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrInvalidAppId     = errors.New("invalid app id")
)

type Auth struct {
	log *zap.Logger
	UserSaver
	UserProvider
	AppProvider
	TokenTTL time.Duration
}

type UserSaver interface {
	Save(ctx context.Context, email string, passwordHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, UserID string) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, AppId int) (models.App, error)
}

// NewAuthService returns new instance of Auth
func NewAuthService(log *zap.Logger, saver UserSaver, provider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		UserSaver:    saver,
		UserProvider: provider,
		AppProvider:  appProvider,
		TokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	user, err := a.UserProvider.User(ctx, email)
	if err != nil {
		if errors.Is(repository.ErrUserNotFound, err) {
			a.log.Warn("User not found", zap.Error(err))
			return "", ErrWrongCredentials
		}
		a.log.Warn("cannot get user", zap.Error(err))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("wrong password", zap.Error(err))
		return "", ErrWrongCredentials
	}

	app, err := a.AppProvider.App(ctx, appID)
	if err != nil {
		a.log.Warn("cannot get app", zap.Error(err))
		return "", err
	}

	a.log.Info("User Logged successfully")

	token, err := jwt.NewJwt(user, app, a.TokenTTL)
	if err != nil {
		a.log.Warn("cannot create token", zap.Error(err))
		return "", err
	}
	return token, nil

}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	a.log.With(zap.String("email", email))
	a.log.Info("Registering new user")

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		a.log.Error("Failed to hash password", zap.Error(err))
		return 0, err
	}

	userId, err := a.UserSaver.Save(ctx, email, bcryptHash)
	if err != nil {
		a.log.Error("Failed to save new user", zap.Error(err))
		return 0, err
	}
	return userId, nil

}

func (a *Auth) IsAdmin(ctx context.Context, UserID string) (bool, error) {
	a.log.Info("Checking admin user", zap.String("UserID", UserID))

	IsAdmin, err := a.UserProvider.IsAdmin(ctx, UserID)
	if err != nil {
		if errors.Is(repository.ErrAppNotFound, err) {
			a.log.Warn("App not found", zap.Error(err))
			return false, ErrInvalidAppId

		}
		a.log.Warn("cannot check admin user", zap.Error(err))
		return false, err
	}
	a.log.Info("Result of checking", zap.Bool("IsAdmin", IsAdmin))

	return IsAdmin, nil

}
