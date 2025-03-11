package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	pb "github.com/daioru/marketplace/internal/generated/api/proto"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	DB *sqlx.DB
}

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

var jwtSecret = []byte("supersecretkey")

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var existingUser User
	err := s.DB.Get(&existingUser, "SELECT * FROM users WHERE email=$1", req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser.ID > 0 {
		return nil, errors.New("email already registered")
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
	err := s.DB.Get(&user, "SELECT * FROM users WHERE email=$1", req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password (db)")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password (bcrypt)")
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	_, err := validateJWT(req.Token)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateTokenResponse{Valid: true}, nil
}

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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
