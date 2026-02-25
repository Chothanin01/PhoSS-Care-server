package entities

import "github.com/google/uuid"

type AdminCreateReq struct {
	Title          string `json:"title"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	CreatedBy      *uuid.UUID `json:"created_by"`
}

type AdminCreateRes struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
}

type AdminUsecase interface {
	CreateAdmin(req *AdminCreateReq) (*AdminCreateRes, error)
}

type AdminRepository interface {
	CreateWithUser(req *AdminCreateReq, hashedPass string) (*AdminCreateRes, error)
}