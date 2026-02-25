package repositories

import (
	"strings"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type DiseaseGetRepository struct {
	db *gorm.DB
}

func NewDiseaseGetRepository(db *gorm.DB) *DiseaseRepository {
	return &DiseaseRepository{db: db}
}

type DiseaseRepository struct {
	db *gorm.DB
}

func NewDiseaseRepository(db *gorm.DB) *DiseaseRepository {
	return &DiseaseRepository{db: db}
}

func (r *DiseaseRepository) GetAllDiseases() ([]entities.Disease, error) {
	var diseases []databases.Disease
	if err := r.db.Find(&diseases).Error; err != nil {
		return nil, err
	}
	var result []entities.Disease
	for _, d := range diseases {
		result = append(result, entities.Disease{
			DiseaseID: d.ID,
			Name:      d.Name,
		})
	}
	return result, nil
}

func (r *DiseaseRepository) LinkPatientDiseases(patientID uuid.UUID, diseases []entities.PatientDiseaseEntity) error {
	var records []databases.PatientDisease
	for _, d := range diseases {
		records = append(records, databases.PatientDisease{
			PatientID: patientID,
			DiseaseID: d.DiseaseID,
			Disease: &databases.Disease{
				Name: d.Name,
			},
		})
	}
	return r.db.Create(&records).Error
}

func (r *PatientGetRepository) GetPatientVaccinesByID(patientID uuid.UUID) (*entities.VaccineInfoRes, error) {
	var vaccines []databases.Vaccine

	// Preload vaccination records and appointments
	if err := r.db.
		Preload("Records.Appoint").
		Find(&vaccines).Error; err != nil {
		return nil, err
	}

	var data []entities.VaccineData

	for _, v := range vaccines {
		var vaccineDetails []entities.VaccineFullInfo

		for _, rec := range v.Records {
			// Filter by patient ID and appointment status
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
			data = append(data, entities.VaccineData{
				VaccineID:   v.ID,
				VaccineName: v.Name,
				Vaccine:     vaccineDetails,
			})
		}
	}

	res := &entities.VaccineInfoRes{
		Success: true,
		Message: "Vaccine info retrieved successfully (only completed appointments)",
		Data:    data,
	}

	return res, nil
}

