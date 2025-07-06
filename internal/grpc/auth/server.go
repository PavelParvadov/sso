package auth

import (
	"context"
	"errors"
	"sso/internal/repository"
	authService "sso/internal/services/auth"

	ssov1 "github.com/PavelParvadov/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int32) (string, error)
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	Auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{Auth: auth})
}

func (s *ServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.Auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, authService.ErrWrongCredentials) {
			return nil, status.Error(codes.InvalidArgument, "wrong credentials")
		}
		return nil, status.Error(codes.Internal, "Internal Error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userId, err := s.Auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		return nil, status.Error(codes.Internal, "Internal Error")
	}
	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}
func (s *ServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	IsAdmin, err := s.Auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Internal Error")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: IsAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email or password is empty")
	}
	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "no app_id provided")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email or password is empty")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == 0 {
		return status.Error(codes.InvalidArgument, "no user_id provided")
	}
	return nil
}
