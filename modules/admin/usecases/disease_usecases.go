package usecases

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type DiseaseUsecase struct {
	repo entities.DiseaseGetRepo
}

func NewDiseaseUsecase(repo entities.DiseaseGetRepo) *DiseaseUsecase {
	return &DiseaseUsecase{repo: repo}
}

func (u *DiseaseUsecase) GetAllDiseases() ([]entities.Disease, error) {
	diseases, err := u.repo.GetAllDiseases()
	if err != nil {
		return nil, fmt.Errorf("get all diseases: %w", err)
	}
	return diseases, nil
}


func (u *patientGetUsecase) GetPatientDiseasesInfo(id uuid.UUID, diseaseID uuid.UUID) (*entities.DiseaseInfoRes, error) {
	patient, err := u.readRepo.GetPatientDiseasesInfoByID(id, diseaseID)
	if err != nil {
		return nil, fmt.Errorf("get patient disease info: %w", err)
	}

	fullname := fmt.Sprintf("%s%s %s", patient.Title, patient.FirstName, patient.LastName)

	var diseaseName string
	for _, pd := range patient.Diseases {
		if pd.Disease != nil && pd.Disease.ID == diseaseID {
			diseaseName = pd.Disease.Name
			break
		}
	}

	healthMap := make(map[uuid.UUID]entities.HealthInfo)
	for _, h := range patient.Healths {
		healthMap[h.AppointID] = entities.HealthInfo{
			Height: float32(h.Height),
			Weight: float32(h.Weight),
			BMI:    h.BMI,
			Pulse:  h.Pulse,
			Sugar:  h.Sugar,
		}
	}
	var diseases []entities.Disease

	for _, pd := range patient.Diseases {
		if pd.Disease != nil {
			diseases = append(diseases, entities.Disease{
				DiseaseID: pd.Disease.ID,
				Name:      pd.Disease.Name,
			})
		}
	}

	var appointments []entities.AppointmentInfo
	for _, ap := range patient.Appointments {
		if ap.DiseaseID == diseaseID && ap.Status == "Completed" {
			var healthPtr *entities.HealthInfo
			if h, ok := healthMap[ap.ID]; ok {
				healthPtr = &h
			}

			appointments = append(appointments, entities.AppointmentInfo{
				No:      ap.No,
				Date:    ap.Date,
				Time:    ap.Time,
				Symptom: ap.Symptom,
				Note:    ap.Note,
				Place:   ap.Place,
				Doctor:  ap.Doctor,
				Status:  ap.Status,
				Letter:  ap.Letter,
				Delay:   ap.Delay,
				Health:  healthPtr,
			})
		}
	}

	data := []entities.DiseaseData{
		{
			PatientID:   patient.ID,
			Fullname:    fullname,
			DiseaseID:   diseaseID,
			DiseaseName: diseaseName,
			Appointment: appointments,
			Disease: diseases,
		},
	}

	res := &entities.DiseaseInfoRes{
		Success: true,
		Message: "Get patient disease info successfully.",
		Data:    data,
	}

	return res, nil
}

func (u *patientGetUsecase) GetPatientVaccinesByID(patientID uuid.UUID) (*entities.VaccineInfoRes, error) {

	patient, vaccines, err := u.readRepo.GetPatientVaccinesByID(patientID)
	if err != nil {
		return nil, fmt.Errorf("get patient vaccine info: %w", err)
	}

	fullname := fmt.Sprintf("%s%s %s", patient.Title, patient.FirstName, patient.LastName)

	var diseases []entities.Disease
	for _, pd := range patient.Diseases {
		if pd.Disease != nil {
			diseases = append(diseases, entities.Disease{
				DiseaseID: pd.Disease.ID,
				Name:      pd.Disease.Name,
			})
		}
	}

	var vaccineList []entities.Vaccine

	for _, v := range vaccines {

		var vaccineDetails []entities.VaccineFullInfo

		for _, rec := range v.Records {

			if rec.Appoint.PatientID != patientID || strings.ToLower(rec.Appoint.Status) != "completed" {
				continue
			}

			vaccineDetails = append(vaccineDetails, entities.VaccineFullInfo{
				RecordID: rec.ID,
				Date:     rec.Appoint.Date.Format("2006-01-02"),
				Type:     v.Type,
				Effect:   v.Effect,
				Note:     v.Note,
				Status:   rec.Status,
				Age:      v.Age,
			})
		}

		if len(vaccineDetails) > 0 {
			vaccineList = append(vaccineList, entities.Vaccine{
				VaccineID:   v.ID,
				VaccineName: v.Name,
				Vaccine:     vaccineDetails,
			})
		}
	}

	data := []entities.VaccineData{
		{
			PatientID: patient.ID,
			Fullname:  fullname,
			Hnnumber:  patient.HnID,
			Vaccine:   vaccineList,
			Disease:   diseases,
		},
	}

	message := "Vaccine info retrieved successfully"
	if len(vaccineList) == 0 {
		message = "No vaccine records found"
	}

	return &entities.VaccineInfoRes{
		Success: true,
		Message: message,
		Data:    data,
	}, nil
}

func (u *DiseaseUsecase) GetAllVaccines() ([]entities.VaccineFullDetail, error) {

	vaccines, err := u.repo.GetAllVaccines()
	if err != nil {
		return nil, fmt.Errorf("get vaccines: %w", err)
	}

	return vaccines, nil
}