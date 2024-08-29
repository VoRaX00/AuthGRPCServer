package auth

import (
	"context"
	ssov1 "github.com/VoRaX00/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appId int32) (token string, err error)
	Register(ctx context.Context, name, email, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
	})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, req.AppId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	validate := validator.New()
	err := validate.Var(req.GetEmail(), "required,email")
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "invalid password")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "invalid app_id")
	}

	return nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.auth.Register(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	validate := validator.New()
	err := validate.Var(req.GetName(), "required,alpha")
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid name")
	}

	err = validate.Var(req.GetEmail(), "required,email")
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "invalid password")
	}

	return nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "invalid user_id")
	}
	return nil
}
