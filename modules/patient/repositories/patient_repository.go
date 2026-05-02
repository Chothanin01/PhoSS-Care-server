package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"gorm.io/gorm"	
)

type GetPatientRepository struct {
	db *gorm.DB
}

func NewGetPatientRepository(db *gorm.DB) *GetPatientRepository {
	return &GetPatientRepository{db: db}
}

func (r *GetPatientRepository) GetPatientBasicInfo(patientID uuid.UUID) (*entities.PatientBasicInfo, error) {
	var dbPatient databases.Patient 

	err := r.db.Select("id, title, first_name, last_name, hn_id, dob").
		Where("id = ?", patientID).
		First(&dbPatient).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound 
		}
		return nil, err
	}

	cleanPatient := &entities.PatientBasicInfo{
		ID:        dbPatient.ID,
		Title:     dbPatient.Title,
		FirstName: dbPatient.FirstName,
		LastName:  dbPatient.LastName,
		HnID:      dbPatient.HnID,
		DOB:       dbPatient.DOB,
	}

	return cleanPatient, nil
}

func (r *GetPatientRepository) GetPatientFullInfo(userID uuid.UUID) (*databases.Patient, error) {
    var patient databases.Patient
    err := r.db.
        Preload("Relatives").
        Preload("Appointments", func(db *gorm.DB) *gorm.DB {
            return db.Order("date DESC, created_at DESC")
        }).
        Where("user_id = ?", userID).
        First(&patient).Error
    
    return &patient, err
}

func (r *GetPatientRepository) GetPatientDiseases(patientID uuid.UUID, filter entities.DiseaseFilter) ([]entities.DiseaseItem, error) {
	var items []entities.DiseaseItem

	query := r.db.Table("patient_disease").
		Select("disease.id as disease_id, disease.name").
		Joins("JOIN disease ON disease.id = patient_disease.disease_id").
		Where("patient_disease.patient_id = ?", patientID)

	if filter.Type == "appoint" {
		query = query.Joins("JOIN appoint ON appoint.disease_id = disease.id AND appoint.patient_id = patient_disease.patient_id")
		
		query = query.Group("disease.id, disease.name")
	}

	err := query.Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}