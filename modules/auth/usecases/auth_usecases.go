package usecases

import (
	"errors"
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/auth/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
)

type authUsecase struct {
	repo     entities.AuthRepository
	passSvc  utils.PasswordService
	jwtSvc   utils.JWTService
	role     string
}

func NewAuthUsecase(repo entities.AuthRepository, passSvc utils.PasswordService, jwtSvc utils.JWTService, role string) entities.AuthUsecase {
	return &authUsecase{repo: repo, passSvc: passSvc, jwtSvc: jwtSvc, role: role}
}

func (u *authUsecase) Login(req *entities.LoginRequest) (*entities.LoginResponse, error) {
    user, err := u.repo.GetByUsername(req.Username)
    if err != nil {
        return nil, errors.New("invalid username or password")
    }

    if !u.passSvc.Compare(user.Password, req.Password) {
        return nil, errors.New("invalid username or password")
    }

    if user.Role != u.role {
        return nil, errors.New("invalid endpoint for this user role")
    }

    roleID, err := u.repo.GetRoleID(user.ID, user.Role)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch role profile: %w", err)
    }

    extraClaims := map[string]interface{}{}

    if user.Role == "patient" {
        diseases, err := u.repo.GetPatientDiseases(user.ID)
        if err != nil {
            return nil, fmt.Errorf("failed to fetch diseases: %w", err)
        }
        extraClaims["diseases"] = diseases
    }

    token, err := u.jwtSvc.GenerateToken(user.ID.String(), user.Role, roleID.String(), extraClaims)
    if err != nil {
        return nil, err
    }

    return &entities.LoginResponse{
        Token:  token,
        Role:   user.Role,
        UserID: user.ID,
    }, nil
}

