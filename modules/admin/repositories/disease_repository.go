package repositories

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiseaseGetRepository struct {
	db *gorm.DB
}

func NewDiseaseGetRepository(db *gorm.DB) *DiseaseGetRepository {
	return &DiseaseGetRepository{db: db}
}

type DiseaseRepository struct {
	db *gorm.DB
}

func NewDiseaseRepository(db *gorm.DB) *DiseaseRepository {
	return &DiseaseRepository{db: db}
}

func (r *DiseaseGetRepository) GetAllDiseases() ([]entities.Disease, error) {
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

func (r *PatientGetRepository) GetPatientVaccinesByID(patientID uuid.UUID) (*databases.Patient, []databases.Vaccine, error) {
	var patient databases.Patient

	err := r.db.
		Preload("Diseases", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("Diseases.Disease", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		First(&patient, "id = ?", patientID).Error
	if err != nil {
		return nil, nil, err
	}

	var vaccines []databases.Vaccine

	err = r.db.
		Preload("Records").
		Preload("Records.Appoint").
		Find(&vaccines).Error
	if err != nil {
		return nil, nil, err
	}

	return &patient, vaccines, nil
}

func (r *DiseaseGetRepository) GetAllVaccines() ([]entities.VaccineFullDetail, error) {

	var records []databases.Vaccine
	var result []entities.VaccineFullDetail

	err := r.db.Find(&records).Error
	if err != nil {
		return nil, err
	}

	for _, rec := range records {

		result = append(result, entities.VaccineFullDetail{
			VaccineID: rec.ID,
			Type:      rec.Type,
			Effect:    rec.Effect,
			Note:      rec.Note,
			Age:       rec.Age,
		})
	}

	return result, nil
}
