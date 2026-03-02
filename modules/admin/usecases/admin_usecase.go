package usecases

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
)

type adminUsecase struct {
	repo     entities.AdminRepository
	passSvc  utils.PasswordService
}

func NewAdminUsecase(repo entities.AdminRepository, passSvc utils.PasswordService) entities.AdminUsecase {
	return &adminUsecase{
		repo:    repo,
		passSvc: passSvc,
	}
}

func (u *adminUsecase) CreateAdmin(req *entities.AdminCreateReq) (*entities.AdminCreateRes, error) {
	if req.Username == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Title == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	hashedPass, err := u.passSvc.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	return u.repo.CreateWithUser(req, hashedPass)
}