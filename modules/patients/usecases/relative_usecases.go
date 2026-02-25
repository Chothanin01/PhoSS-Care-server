package usecases

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

type RelativeUpdateUsecase interface {
	UpdateAllRelatives(patientID uuid.UUID, req *entities.RelativeAllUpdateReq) (*entities.RelativeAllUpdateRes, error)
}

type relativeUpdateUsecase struct {
	repo entities.RelativeUpdateRepo
}

func NewRelativeUsecase(repo entities.RelativeUpdateRepo) entities.RelativeUsecase {
	return &relativeUpdateUsecase{repo: repo}
}

func (u *relativeUpdateUsecase) UpdateAllRelatives(patientID uuid.UUID, req *entities.RelativeAllUpdateReq) (*entities.RelativeAllUpdateRes, error) {
	if req.UpdatedBy == uuid.Nil {
		return nil, fmt.Errorf("missing updated_by field")
	}

	if req.Kin.FirstName == "" || req.Kin.LastName == "" ||
		req.Caretaker.FirstName == "" || req.Caretaker.LastName == "" ||
		req.Medicine.FirstName == "" || req.Medicine.LastName == "" {
		return nil, fmt.Errorf("all relative sections (kin, caretaker, medicine) must be fully provided")
	}

	fields := map[string]entities.RelativeUpdateReq{
		"kin":       req.Kin,
		"caretaker": req.Caretaker,
		"medicine":  req.Medicine,
	}
	for role, r := range fields {
		if r.PhoneNumber == "" ||
			r.Address.HouseNumber == "" ||
			r.Address.SubDistrict == "" ||
			r.Address.District == "" ||
			r.Address.Province == "" ||
			r.Address.ZipCode == "" {
			return nil, fmt.Errorf("missing field(s) in %s section", role)
		}
	}

	return u.repo.UpdateRelativeInfo(patientID, req)
}

func (u *relativeUpdateUsecase) UpdateAllOfficers(patientID uuid.UUID, req *entities.OfficerAllUpdateReq) (*entities.OfficerAllUpdateRes, error) {
	if req.UpdatedBy == uuid.Nil {
		return nil, fmt.Errorf("missing updated_by field")
	}

	if req.House.FirstName == "" || req.House.LastName == "" ||
		req.Nurse.FirstName == "" || req.Nurse.LastName == "" {
		return nil, fmt.Errorf("both officer sections (house, nurse) must be fully provided")
	}

	if req.House.Title == "" || req.Nurse.Title == "" {
		return nil, fmt.Errorf("title is required for both house and nurse officers")
	}

	return u.repo.UpdateOfficerInfo(patientID, req)
}