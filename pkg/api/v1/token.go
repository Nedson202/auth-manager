package v1

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JwtClaims struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expirationTime"`
	jwt.StandardClaims
}

type TokenService interface {
	GenerateToken(user *UserWithoutPassword) (string, error)
	VerifyToken(token string) (jwt.Claims, error)
	GetAuthorizationToken(ctx context.Context) (string, error)
}

type tokenService struct {
	config Config
}

func NewTokenService(config Config) *tokenService {
	return &tokenService{config}
}

func (r *tokenService) GenerateToken(user *UserWithoutPassword) (string, error) {
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := JwtClaims{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(r.config.JwtSecret))
	if err != nil {
		log.Println("token-service:GenerateToken->" + err.Error())
	}
	return token, err
}

func (r *tokenService) VerifyToken(token string) (jwt.Claims, error) {
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "Invalid authorization header-> ")
	}

	decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("an error occured verifying token")
		}
		return r.config.JwtSecret, nil
	})
	if err != nil {
		log.Println("token-service:VerifyToken->" + err.Error())
		return nil, status.Error(codes.Unauthenticated, "invalid authorization token-> "+err.Error())
	}

	if !decodedToken.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
	}

	return decodedToken.Claims, nil
}

func (r *tokenService) GetAuthorizationToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	return values[0], nil
}
