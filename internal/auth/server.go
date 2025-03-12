package auth

import (
	"context"
	"time"

	pb "github.com/daioru/marketplace/internal/generated/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	DB *sqlx.DB
}

type User struct {
	ID             int        `db:"id"`
	Email          string     `db:"email"`
	PasswordHash   string     `db:"password_hash"`
	RefreshToken   *string    `db:"refresh_token"`      // Может быть NULL
	RefreshExpires *time.Time `db:"refresh_expires_at"` // Может быть NULL
	CreatedAt      time.Time  `db:"created_at"`
}

var jwtSecret = []byte("supersecretkey")

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var exists bool
	err := s.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing user")
	}
	if exists {
		return nil, status.Errorf(codes.AlreadyExists, "email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var resp pb.RegisterResponse
	err = s.DB.QueryRow("INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		req.Username, req.Email, string(hashedPassword)).Scan(&resp.UserId)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user User
	err := s.DB.Get(&user, "SELECT id, password_hash FROM users WHERE email=$1", req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	accessToken, refreshToken, err := GenerateTokens(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens")
	}

	_, err = s.DB.Exec("UPDATE users SET refresh_token=$1 WHERE id=$2", refreshToken, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save refresh token")
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	_, err := validateJWT(req.Token)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateTokenResponse{Valid: true}, nil
}

func validateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	var userID int
	var refreshExpiresAt time.Time

	err := s.DB.QueryRow("SELECT id, refresh_expires_at FROM users WHERE refresh_token=$1", req.RefreshToken).
		Scan(&userID, &refreshExpiresAt)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	if time.Now().After(refreshExpiresAt) {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token expired")
	}

	accessToken, newRefreshToken, err := GenerateTokens(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate new tokens")
	}

	newRefreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.DB.Exec("UPDATE users SET refresh_token=$1, refresh_expires_at=$2 WHERE id=$3",
		newRefreshToken, newRefreshExpiresAt, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	_, err := s.DB.Exec("UPDATE users SET refresh_token=NULL, refresh_expires_at=NULL WHERE refresh_token=$1", req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout")
	}

	return &pb.LogoutResponse{Message: "Logout successful"}, nil
}
