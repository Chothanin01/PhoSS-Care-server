package entities

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Password string
	Role     string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	UserID uuid.UUID `json:"user_id"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}


type AuthRepository interface {
	GetByUsername(username string) (*User, error)
	GetPatientDiseases(userID uuid.UUID) ([]map[string]interface{}, error)
	GetRoleID(userID uuid.UUID, role string) (uuid.UUID, error)
}

type AuthUsecase interface {
	Login(req *LoginRequest) (*LoginResponse, error)
}

