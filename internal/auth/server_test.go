package auth_test

import (
	"context"
	"log"
	"testing"

	"github.com/daioru/marketplace/internal/auth"
	pb "github.com/daioru/marketplace/internal/generated/api/proto"
	migrations "github.com/daioru/marketplace/migrations/auth"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	DB     *sqlx.DB
	Server *auth.AuthService
}

func (s *AuthTestSuite) SetupSuite() {
	dsn := "postgres://auth_test_user:auth_test_pass@localhost:5433/auth_test_db?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connetction error: %v", err)
	}
	s.DB = db

	goose.SetBaseFS(migrations.AuthEmbedFS)

	err = goose.Up(db.DB, ".")
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	s.Server = &auth.AuthService{DB: db}
}

func (s *AuthTestSuite) TearDownSuite() {
	s.DB.Close()
}

func (s *AuthTestSuite) TestRegisterAndLogin() {
	s.T().Run("RegisterUser", func(t *testing.T) {
		req := &pb.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		resp, err := s.Server.Register(context.Background(), req)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.UserId)

		var count int
		err = s.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1", req.Email)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1, count)
	})

	s.T().Run("LoginUser", func(t *testing.T) {
		req := &pb.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		resp, err := s.Server.Login(context.Background(), req)

		assert.NoError(s.T(), err)
		assert.NotEmpty(s.T(), resp.AccessToken)
		assert.NotEmpty(s.T(), resp.AccessToken)
	})
}



func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
