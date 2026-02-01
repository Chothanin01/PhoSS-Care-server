package usecases

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

type patientUsecase struct {
	repo entities.PatientRepository
}

func NewPatientUsecase(repo entities.PatientRepository) entities.PatientUsecase {
	return &patientUsecase{
		repo: repo,
	}
}

func (u *patientUsecase) Create(req *entities.PatientCreateReq) (*entities.PatientCreateRes, error) {
	return u.repo.Create(req)
}
