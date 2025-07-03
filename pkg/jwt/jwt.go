package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/domain/models"
	"time"
)

func NewJwt(user models.User, app models.App, t time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(t).Unix()
	claims["email"] = user.Email
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
