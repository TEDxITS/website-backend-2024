package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userId string, role string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetPayloadInsideToken(token string) (string, string, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "TEDxITS 2024",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(userId string, role string) string {
	claims := jwtCustomClaim{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * constants.JWT_EXPIRE_TIME_IN_MINUTES)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return tx
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) GetPayloadInsideToken(token string) (string, string, error) {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return "", "", err
	}

	if !t_Token.Valid {
		return "", "", dto.ErrTokenInvalid
	}

	claims := t_Token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	role := fmt.Sprintf("%v", claims["role"])
	return id, role, nil
}
