package usecases

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

type PasswordService interface {
	Hash(password string) (string, error)
}

type patientUsecase struct {
	tx          entities.Transaction
	passwordSvc PasswordService
}

func NewPatientUsecase(tx entities.Transaction, passwordSvc PasswordService) *patientUsecase {
	return &patientUsecase{
		tx:          tx,
		passwordSvc: passwordSvc,
	}
}

func (u *patientUsecase) CreateFull(req *entities.PatientFullCreateReq) (*entities.PatientCreateRes, error) {
	var res *entities.PatientCreateRes

	err := u.tx.Do(func(r entities.RepositorySet) error {
		nextHnID, err := r.PatientRepo.GenerateNextHnID()
		if err != nil {
			return err
		}

		hnidStr := fmt.Sprintf("%s", nextHnID)
		hashedPass, err := u.passwordSvc.Hash(hnidStr)
		if err != nil {
			return fmt.Errorf("failed to hash hnid: %w", err)
		}

		user, err := r.UserRepo.Create(req.Patient.IDCard, hashedPass, "patient")
		if err != nil {
			return err
		}

		patientReq := &entities.PatientCreateReq{
			Title:       req.Patient.Title,
			FirstName:   req.Patient.FirstName,
			LastName:    req.Patient.LastName,
			Dob:         req.Patient.DOB,
			HnID:        nextHnID,
			IDCard:      req.Patient.IDCard,
			Nationality: req.Patient.Nationality,
			Ethnicity:   req.Patient.Ethnicity,
			PhoneNumber: req.Patient.PhoneNumber,
			Address:     req.Patient.Address,
			Allergy:     req.Patient.Allergy,
			Rights:      req.Patient.Rights,
			UserID:      user.ID,
			CreatedBy:   req.CreatedBy, 
			UpdatedBy:   req.CreatedBy,
		}

		res, err = r.PatientRepo.CreateWithUser(patientReq, user.ID)
		if err != nil {
			return err
		}
	
		if len(req.Patient.Diseases) > 0 {
			var diseases []entities.PatientDiseaseEntity
			for _, d := range req.Patient.Diseases {
				diseases = append(diseases, entities.PatientDiseaseEntity{
					DiseaseID: d.DiseaseID,
					Name:      d.Name,
				})
			}
			if err := r.DiseaseRepo.LinkPatientDiseases(uint(res.Id), diseases); err != nil {
				return err
			}
		}

		relatives := []entities.RelativeEntity{
			makeRelative(req.Relative.Kin, "kin", res.Id, req.CreatedBy),
			makeRelative(req.Relative.Caretaker, "caretaker", res.Id, req.CreatedBy),
			makeRelative(req.Relative.Medicine, "medicine", res.Id, req.CreatedBy),
			makeRelative(req.Officer.House, "house", res.Id, req.CreatedBy),
			makeRelative(req.Officer.Nurse, "nurse", res.Id, req.CreatedBy),
		}


		return r.RelativeRepo.Create(relatives)
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}

func makeRelative(d entities.RelativeDetail, role string, pid uint64, creator uint) entities.RelativeEntity {
	return entities.RelativeEntity{
		Title:       d.Title,
		FirstName:   d.FirstName,
		LastName:    d.LastName,
		PhoneNumber: d.PhoneNumber,
		Role:        role,
		PatientID:   uint(pid),
		Address: entities.AddressReq{
			HouseNumber:  d.Address.HouseNumber,
			VillageNumber: d.Address.VillageNumber,
			Alley:         d.Address.Alley,
			Road:         d.Address.Road,
			SubDistrict:  d.Address.SubDistrict,
			District:     d.Address.District,
			Province:     d.Address.Province,
			ZipCode:      d.Address.ZipCode,
		},
		CreatedBy: creator,
		UpdatedBy: creator,
	}
}
