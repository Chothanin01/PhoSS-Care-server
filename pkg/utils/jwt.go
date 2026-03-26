package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID string, role string, role_id string, extraClaims map[string]interface{}) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetSecret() string
}

type jwtService struct {
	secret string
	issuer string
}

func NewJWTService(secret, issuer string) JWTService {
	return &jwtService{
		secret: secret,
		issuer: issuer,
	}
}

func (j *jwtService) GenerateToken(userID string, role string, role_id string, extraClaims map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": role_id,
		"role":    role,
		"iss":     j.issuer,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	for k, v := range extraClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
}

func (j *jwtService) GetSecret() string {
	return j.secret
}