package usecases

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

// ---------------------- CREATE ----------------------

type PasswordService interface {
	Hash(password string) (string, error)
}

type newPatientUsecase struct {
	tx          entities.Transaction
	passwordSvc PasswordService
}

func NewPatientUsecase(tx entities.Transaction, passwordSvc PasswordService) *newPatientUsecase {
	return &newPatientUsecase{
		tx:          tx,
		passwordSvc: passwordSvc,
	}
}

func (u *newPatientUsecase) CreateFull(req *entities.PatientFullCreateReq) (*entities.PatientCreateRes, error) {
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
			HouseNumber:   d.Address.HouseNumber,
			VillageNumber: d.Address.VillageNumber,
			Alley:         d.Address.Alley,
			Road:          d.Address.Road,
			SubDistrict:   d.Address.SubDistrict,
			District:      d.Address.District,
			Province:      d.Address.Province,
			ZipCode:       d.Address.ZipCode,
		},
		CreatedBy: creator,
		UpdatedBy: creator,
	}
}

// ---------------------- GET ----------------------

type patientGetUsecase struct {
	readRepo entities.PatientGetRepo
}

func NewPatientGetUsecase(readRepo entities.PatientGetRepo) entities.PatientGetUsecase {
	return &patientGetUsecase{readRepo: readRepo}
}

func (u *patientGetUsecase) GetPatientList(page, limit int) (*entities.PatientListRes, error) {
	params := entities.PatientQueryParams{
		Page:  page,
		Limit: limit,
	}
	return u.GetPatientListWithFilter(params)
}

func (u *patientGetUsecase) GetPatientListWithFilter(req entities.PatientQueryParams) (*entities.PatientListRes, error) {
	dbPatients, err := u.readRepo.GetPatientsWithFilter(req)
	if err != nil {
		return nil, err
	}

	var total int64
	if req.Search != "" || len(req.Diseases) > 0 || req.Appoint != nil {
		total, err = u.readRepo.CountPatientsWithFilter(req)
	} else {
		total, err = u.readRepo.CountPatients()
	}
	if err != nil {
		return nil, err
	}

	res := &entities.PatientListRes{
		Success:    true,
		Message:    "Patient list fetched successfully",
		Page:       req.Page,
		PerPage:    req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
		Data:       []entities.PatientInfo{},
	}

	for _, p := range dbPatients {
		pInfo := entities.PatientInfo{
			ID:       p.ID,
			FullName: p.Title + "" + p.FirstName + " " + p.LastName,
			IDCard:   p.IDCard,
			HnNumber: p.HnID,
		}

		hasAnyAppointment := false
		for _, pd := range p.Diseases {
			hasOngoing := false
			for _, ap := range p.Appointments {
				if ap.DiseaseID == pd.DiseaseID && ap.Status == "ongoing" {
					hasOngoing = true
					hasAnyAppointment = true
					break
				}
			}
			pInfo.Diseases = append(pInfo.Diseases, entities.DiseaseWithStatus{
				DiseaseID:      pd.DiseaseID,
				Name:           pd.Disease.Name,
				HasAppointment: hasOngoing,
			})
		}

		if req.Appoint != nil {
			if *req.Appoint {
				if !hasAnyAppointment {
					continue
		}
		} else {
			if hasAnyAppointment {
				continue
		}
	}
}
		res.Data = append(res.Data, pInfo)
	}

	return res, nil
}
