package services

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"challenge.go.lgsjesus/application/dtos"
	"challenge.go.lgsjesus/domain"
	"challenge.go.lgsjesus/framework/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtKey []byte
var jwtIssuer string
var jwtAudience string
var jwtExpirationTime int

type AuthService struct {
	repository *repositories.UserRepositoryDb
}

func NewAuthService(repository *repositories.UserRepositoryDb) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the JWT key or any other necessary setup
	// This can be moved to a config file or environment variable in production
	jwtKey = []byte(os.Getenv("JWT_SECRET")) // Replace with your secret key
	jwtIssuer = os.Getenv("JWT_ISSUER")
	jwtAudience = os.Getenv("JWT_AUDIENCE")
	if jwtKey == nil || jwtIssuer == "" || jwtAudience == "" {
		log.Fatal("JWT configuration is not set properly. Please check your environment variables.")
	}
	stringExpiration := os.Getenv("JWT_EXPIRATION_TIME")
	if stringExpiration == "" {
		jwtExpirationTime = 10 // Default to 10 min if not set
	} else {
		var err error
		jwtExpirationTime, err = strconv.Atoi(stringExpiration)
		if err != nil {
			log.Fatalf("Invalid JWT_EXPIRATION_TIME: %v", err)
		}
	}
}
func (s *AuthService) AuthenticateUser(authDto *dtos.AuthDto) (*dtos.TokenDto, error) {
	user, err := s.repository.FindByNickName(authDto.NickName)
	if err != nil {
		return nil, err
	}

	if user == nil || !checkPassword(user.Password, authDto.Password) {
		return nil, errors.New("invalid credentials")
	}

	return generateJWT(user)
}

func generateJWT(user *domain.User) (*dtos.TokenDto, error) {

	expiration := jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpirationTime) * time.Minute))

	claims := &jwt.RegisteredClaims{
		Subject:   user.NickName,
		Issuer:    jwtIssuer,
		Audience:  jwt.ClaimStrings{jwtAudience},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenAssigned, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &dtos.TokenDto{
		Token:     tokenAssigned,
		ExpiresAt: expiration.Unix(),
	}, nil
}
func GetJWTSecret() string {
	if jwtKey == nil {
		log.Fatal("JWT secret key is not set. Please check your environment variables.")
	}
	return string(jwtKey)
}
