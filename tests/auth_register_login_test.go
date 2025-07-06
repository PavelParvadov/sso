package tests

import (
	"sso/tests/suite"
	"testing"
	"time"

	ssov1 "github.com/PavelParvadov/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppId        = 0
	appId             = 1
	defaultPassLength = 10
	secret            = "pavel"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, s := suite.NewSuite(t)
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, defaultPassLength)
	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	loginResp, err := s.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})
	require.NoError(t, err)
	tm := time.Now()

	token := loginResp.GetToken()

	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))
	assert.Equal(t, regResp.GetUserId(), int64(claims["uid"].(float64)))

	const deltaSecond = 1

	assert.InDelta(t, tm.Add(s.Cfg.TokenTTL).Unix(), claims["exp"], deltaSecond)

}

func TestRegisterLogin_DuplicateRegister(t *testing.T) {
	ctx, s := suite.NewSuite(t)
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, defaultPassLength)
	regResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.GetUserId())

	DuplicateRegResp, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, DuplicateRegResp.GetUserId())

}

func TestRegister_FailCases(t *testing.T) {
	ctx, s := suite.NewSuite(t)
	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "email or password is empty",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, true, defaultPassLength),
			expectedErr: "email or password is empty",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "email or password is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, s := suite.NewSuite(t)
	tests := []struct {
		name        string
		email       string
		password    string
		AppId       int32
		ExpectedErr string
	}{
		{
			name:        "Test with incorrect password",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, true, defaultPassLength),
			AppId:       appId,
			ExpectedErr: "Internal Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			registeredEmail := gofakeit.Email()
			_, err := s.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    registeredEmail,
				Password: gofakeit.Password(true, true, true, true, true, defaultPassLength),
			})
			require.NoError(t, err)

			_, err = s.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    registeredEmail,
				Password: tt.password,
				AppId:    tt.AppId,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.ExpectedErr)
		})
	}
}
