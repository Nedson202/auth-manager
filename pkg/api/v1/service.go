package v1

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
	v1 "github.com/nedson202/auth-manager/api/proto/v1"
	"github.com/nedson202/auth-manager/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// authServiceServer is implementation of v1.AuthServiceServer proto interface
type authServiceServer struct {
	repository   Repository
	tokenService TokenService
	v1.UnimplementedAuthServiceServer
}

// NewAuthServiceServer creates Todo service
func NewAuthServiceServer(repo Repository, tokenService TokenService) v1.AuthServiceServer {
	return &authServiceServer{repo, tokenService, v1.UnimplementedAuthServiceServer{}}
}

// Create new todo task
func (s *authServiceServer) Signup(ctx context.Context, req *v1.AuthRequest) (*v1.AuthResponse, error) {
	hashedPassword, _ := hashPassword(req.Password)
	userID := getUUID()

	user, err := s.repository.CreateUser(&v1.User{
		Id:       userID,
		Username: "",
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		logger.Log.Error("service:Signup:::failed to create user-> " + err.Error())
		return nil, status.Error(codes.Unknown, "failed to create user")
	}

	if user.Id == "" {
		return nil, status.Error(codes.AlreadyExists, "please provide another email")
	}

	token, _ := s.tokenService.GenerateToken(&UserWithoutPassword{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	})
	return &v1.AuthResponse{
		Id:    user.GetId(),
		Token: token,
	}, nil
}

// Login user
func (s *authServiceServer) Login(ctx context.Context, req *v1.AuthRequest) (*v1.AuthResponse, error) {
	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		logger.Log.Error("service:Login:::error during user login-> " + err.Error())
		return nil, status.Error(codes.Unauthenticated, "authentication failed due to invalid credentials")
	}
	if user.Id == "" {
		return nil, status.Error(codes.Unauthenticated, "authentication failed due to invalid credentials")
	}

	isPassword := verifyPassword(req.Password, user.Password)
	if !isPassword {
		return nil, status.Error(codes.Unauthenticated, "authentication failed due to invalid credentials")
	}

	token, _ := s.tokenService.GenerateToken(&UserWithoutPassword{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	})
	return &v1.AuthResponse{
		Id:    user.GetId(),
		Token: token,
	}, nil
}

func (s *authServiceServer) RefreshToken(ctx context.Context, req *v1.TokenRefreshRequest) (*v1.AuthResponse, error) {
	user := &v1.User{}
	token, err := s.tokenService.GetAuthorizationToken(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	refreshedTokenClaims, err := s.tokenService.VerifyToken(token)
	mapstructure.Decode(refreshedTokenClaims, &user)

	if err != nil {
		logger.Log.Error("service:RefreshToken:::failed to refresh token-> " + err.Error())
		return nil, status.Error(codes.PermissionDenied, "Token refresh failed")
	}

	refreshedToken, _ := s.tokenService.GenerateToken(&UserWithoutPassword{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	})
	return &v1.AuthResponse{
		Id:    user.Id,
		Token: refreshedToken,
	}, nil
}
