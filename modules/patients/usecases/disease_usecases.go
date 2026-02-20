package usecases

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
)

type newDiseaseUsecase struct {
	repo entities.DiseaseGetRepo
}

func NewDiseaseUsecase(repo entities.DiseaseGetRepo) *newDiseaseUsecase {
	return &newDiseaseUsecase{repo: repo}
}

func (u *newDiseaseUsecase) GetAllDiseases() ([]entities.Disease, error) {
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
		},
	}

	res := &entities.DiseaseInfoRes{
		Success: true,
		Message: "Get patient disease info successfully.",
		Data:    data,
	}

	return res, nil
}