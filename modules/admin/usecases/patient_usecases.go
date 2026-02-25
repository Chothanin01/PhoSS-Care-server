package usecases

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
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

func (u *newPatientUsecase) CreateFull(req *entities.PatientFullCreateReq, creatorID uuid.UUID) (*entities.PatientCreateRes, error) {
	var res *entities.PatientCreateRes

	p := req.Patient
	if p.FirstName == "" || p.LastName == "" ||
		p.Sex == "" || p.Title == "" ||
		p.DOB == "" || p.IDCard == "" ||
		p.Nationality == "" || p.Ethnicity == "" ||
		p.Rights == "" || p.Weight <= 0 || p.Height <= 0 ||
		p.PhoneNumber == "" ||
		p.Address.HouseNumber == "" || p.Address.SubDistrict == "" ||
		p.Address.District == "" || p.Address.Province == "" ||
		p.Address.ZipCode == "" || len(p.Diseases) == 0 ||
		(req.Relative.Kin.FirstName == "" && req.Relative.Kin.LastName == "") {
		return nil, fmt.Errorf("missing required patient information")
	}

	r := req.Relative.Kin
	if r.FirstName == "" || r.LastName == "" ||
		r.PhoneNumber == "" ||
		r.Address.HouseNumber == "" || r.Address.SubDistrict == "" ||
		r.Address.District == "" || r.Address.Province == "" || r.Address.ZipCode == "" {
		return nil, fmt.Errorf("missing required kin information")
	}

	err := u.tx.Do(func(r entities.RepositorySet) error {
		nextHnID, err := r.PatientRepo.GenerateNextHnID()
		if err != nil {
			return fmt.Errorf("generate HN ID: %w", err)
		}

		hashedPass, err := u.passwordSvc.Hash(nextHnID)
		if err != nil {
			return fmt.Errorf("hash HNID: %w", err)
		}

		user, err := r.UserRepo.Create(p.IDCard, hashedPass, "patient", &creatorID)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		patientReq := &entities.PatientCreateReq{
			Title:       p.Title,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			Sex:         p.Sex,
			Dob:         p.DOB,
			HnID:        nextHnID,
			IDCard:      p.IDCard,
			Nationality: p.Nationality,
			Ethnicity:   p.Ethnicity,
			PhoneNumber: p.PhoneNumber,
			Address:     p.Address,
			Allergy:     p.Allergy,
			Rights:      p.Rights,
			Weight:      p.Weight,
			Height:      p.Height,
			UserID:      user.ID,
		}

		res, err = r.PatientRepo.CreateWithUser(patientReq, user.ID, creatorID)
		if err != nil {
			return fmt.Errorf("create patient: %w", err)
		}

		if len(p.Diseases) > 0 {
			diseases := make([]entities.PatientDiseaseEntity, 0, len(p.Diseases))
			for _, d := range p.Diseases {
				diseases = append(diseases, entities.PatientDiseaseEntity{
					DiseaseID: d.DiseaseID,
					Name:      d.Name,
				})
			}
			if err := r.DiseaseRepo.LinkPatientDiseases(res.Id, diseases); err != nil {
				return fmt.Errorf("link diseases: %w", err)
			}
		}

		relatives := []entities.RelativeEntity{
			makeRelative(req.Relative.Kin, "kin", res.Id, creatorID),
			makeRelative(req.Relative.Caretaker, "caretaker", res.Id, creatorID),
			makeRelative(req.Relative.Medicine, "medicine", res.Id, creatorID),
			makeRelative(req.Officer.House, "house", res.Id, creatorID),
			makeRelative(req.Officer.Nurse, "nurse", res.Id, creatorID),
		}

		if err := r.RelativeRepo.Create(relatives); err != nil {
			return fmt.Errorf("create relatives: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}


func makeRelative(d entities.RelativeCreate, role string, pid uuid.UUID, creator uuid.UUID) entities.RelativeEntity {
	return entities.RelativeEntity{
		Title:       d.Title,
		FirstName:   d.FirstName,
		LastName:    d.LastName,
		PhoneNumber: d.PhoneNumber,
		Role:        role,
		PatientID:   pid,
		Address: entities.Address{
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
		Data:       []entities.PatientHomeInfo{},
	}

	for _, p := range dbPatients {
		pInfo := entities.PatientHomeInfo{
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

func (u *patientGetUsecase) GetPatientInfoByID(id uuid.UUID) (*entities.PatientInfoRes, error) {

	patient, err := u.readRepo.GetPatientInfoByID(id)
	if err != nil {
		return nil, fmt.Errorf("get patient by id: %w", err)
	}

	addr := entities.Address{
		HouseNumber:   patient.Address.HouseNumber,
		VillageNumber: patient.Address.VillageNumber,
		Alley:         patient.Address.Alley,
		Road:          patient.Address.Road,
		SubDistrict:   patient.Address.SubDistrict,
		District:      patient.Address.District,
		Province:      patient.Address.Province,
		ZipCode:       patient.Address.ZipCode,
	}

	var ageYears, ageMonths, ageDays int

	if !patient.DOB.IsZero() {
		now := time.Now()
		years := now.Year() - patient.DOB.Year()
		months := int(now.Month()) - int(patient.DOB.Month())
		days := now.Day() - patient.DOB.Day()

		if days < 0 {
			prevMonth := now.AddDate(0, -1, 0)
			days += utils.DaysInMonth(prevMonth.Year(), prevMonth.Month())
			months--
		}
		if months < 0 {
			months += 12
			years--
		}

		ageYears = years
		ageMonths = months
		ageDays = days
	}

	full := entities.PatientFullInfo{
		Fullname:    fmt.Sprintf("%s%s %s", patient.Title, patient.FirstName, patient.LastName),
		Sex:         patient.Sex,
		IDCard:      patient.IDCard,
		HnNumber:    patient.HnID,
		Rights:      patient.Rights,
		AgeYears:    ageYears,
		AgeMonths:   ageMonths,
		AgeDays:     ageDays,
		Allergy:     patient.Allergy,
		PhoneNumber: patient.PhoneNumber,
		Address:     addr,
		Weight:      patient.Weight,
		Height:      patient.Height,
	}

	var diseases []entities.Disease
	for _, pd := range patient.Diseases {
		diseases = append(diseases, entities.Disease{
			DiseaseID:   pd.DiseaseID,
			Name: pd.Disease.Name,
		})
	}

	var kin, caretaker, medicine entities.RelativeInfo
	var house, nurse entities.OfficerInfo

	for _, rel := range patient.Relatives {
	fullName := fmt.Sprintf("%s%s %s", rel.Title, rel.FirstName, rel.LastName)

	switch rel.Role {
		case "kin":
			kin = entities.RelativeInfo{
				Fullname:    fullName,
				PhoneNumber: rel.PhoneNumber,
				Role:        rel.Role,
				Address: entities.Address{
					HouseNumber:   rel.Address.HouseNumber,
					VillageNumber: rel.Address.VillageNumber,
					Alley:         rel.Address.Alley,
					Road:          rel.Address.Road,
					SubDistrict:   rel.Address.SubDistrict,
					District:      rel.Address.District,
					Province:      rel.Address.Province,
					ZipCode:       rel.Address.ZipCode,
				},
			}
		case "caretaker":
			caretaker = entities.RelativeInfo{
				Fullname:    fullName,
				PhoneNumber: rel.PhoneNumber,
				Role:        rel.Role,
				Address: entities.Address{
					HouseNumber:   rel.Address.HouseNumber,
					VillageNumber: rel.Address.VillageNumber,
					Alley:         rel.Address.Alley,
					Road:          rel.Address.Road,
					SubDistrict:   rel.Address.SubDistrict,
					District:      rel.Address.District,
					Province:      rel.Address.Province,
					ZipCode:       rel.Address.ZipCode,
				},
			}
		case "medicine":
			medicine = entities.RelativeInfo{
				Fullname:    fullName,
				PhoneNumber: rel.PhoneNumber,
				Role:        rel.Role,
				Address: entities.Address{
					HouseNumber:   rel.Address.HouseNumber,
					VillageNumber: rel.Address.VillageNumber,
					Alley:         rel.Address.Alley,
					Road:          rel.Address.Road,
					SubDistrict:   rel.Address.SubDistrict,
					District:      rel.Address.District,
					Province:      rel.Address.Province,
					ZipCode:       rel.Address.ZipCode,
				},
			}
		case "house":
			house = entities.OfficerInfo{
				Fullname: fullName,
				Role:     rel.Role,
			}
		case "nurse":
			nurse = entities.OfficerInfo{
				Fullname: fullName,
				Role:     rel.Role,
			}
		}
	}


	res := &entities.PatientInfoRes{
		Success: true,
		Message: "Get patient info successfully.",
		Data: []entities.PatientData{
			{
				Patient:  full,
				Disease:  diseases,
				Relative: entities.Relative{
					Kin:       kin,
					Caretaker: caretaker,
					Medicine:  medicine,
				},
				Officer: entities.Officer{
					House: house,
					Nurse: nurse,
				},
			},
		},
	}

	return res, nil
}

func (u *patientGetUsecase) GetPatientAppointmentsByID(id uuid.UUID) (*entities.AppointInfoRes, error) {
	patient, err := u.readRepo.GetPatientAppointmentsByID(id)
	if err != nil {
		return nil, fmt.Errorf("get patient appointments: %w", err)
	}

	fullname := fmt.Sprintf("%s%s %s", patient.Title, patient.FirstName, patient.LastName)

	diseaseMap := make(map[uuid.UUID]*entities.AppointDisease)
	for _, ap := range patient.Appointments {
		if ap.Status != "Ongoing" && ap.Status != "Delay" {
			continue
		}

		d, ok := diseaseMap[ap.DiseaseID]
		if !ok {
			diseaseMap[ap.DiseaseID] = &entities.AppointDisease{
				DiseaseID:   ap.Disease.ID,
				DiseaseName: ap.Disease.Name,
				Appointments: []entities.AppointmentFullInfo{},
			}
			d = diseaseMap[ap.DiseaseID]
		}

		d.Appointments = append(d.Appointments, entities.AppointmentFullInfo{
			No:      ap.No,
			Date:    ap.Date.Format("2006-01-02"),
			Time:    ap.Time,
			Symptom: ap.Symptom,
			Note:    ap.Note,
			Place:   ap.Place,
			Doctor:  ap.Doctor,
			Status:  ap.Status,
			Letter:  ap.Letter,
			Delay:   ap.Delay,
		})
	}

	var diseases []entities.AppointDisease
	for _, d := range diseaseMap {
		diseases = append(diseases, *d)
	}

	res := &entities.AppointInfoRes{
		Success: true,
		Message: "Get patient appointments successfully.",
		Data: []entities.AppointData{
			{
				PatientID: patient.ID,
				Fullname:  fullname,
				Hnnumber:  patient.HnID,
				Diseases:  diseases,
			},
		},
	}

	return res, nil
}

// ---------------------- EDIT ----------------------
func (u *newPatientUsecase) UpdatePatientInfo(id uuid.UUID, req *entities.PatientUpdateReq, adminID uuid.UUID) (*entities.PatientUpdateRes, error) {

	if req.FirstName == "" || req.LastName == "" ||
		req.Sex == "" || req.Title == "" ||
		req.DOB == "" || req.IDCard == "" ||
		req.Rights == "" || req.Nationality == "" || req.Ethnicity == "" ||
		req.PhoneNumber == "" ||
		req.Address.HouseNumber == "" || req.Address.SubDistrict == "" ||
		req.Address.District == "" || req.Address.Province == "" || req.Address.ZipCode == "" {
		return nil, fmt.Errorf("missing required patient information")
	}

	var res *entities.PatientUpdateRes

	err := u.tx.Do(func(r entities.RepositorySet) error {
		updated, err := r.PateintUpdateRepo.UpdatePatientInfo(id, req, adminID)
		if err != nil {
			return fmt.Errorf("update patient info: %w", err)
		}
		res = updated
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
