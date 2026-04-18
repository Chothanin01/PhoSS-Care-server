package repositories

import (
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"	
)

type GetPatientRepository struct {
	db *gorm.DB
}

func NewGetPatientRepository(db *gorm.DB) *GetPatientRepository {
	return &GetPatientRepository{db: db}
}

func (r *GetPatientRepository) GetPatientBasicInfo(patientID uuid.UUID) (*databases.Patient, error) {
	var patient databases.Patient
	
	err := r.db.
		First(&patient, "id = ?", patientID).Error

	if err != nil {
		return nil, err
	}

	return &patient, nil
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