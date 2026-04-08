package repositories

import (
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"	
)

type PatientGetRepository struct {
	db *gorm.DB
}

func NewPatientGetRepository(db *gorm.DB) *PatientGetRepository {
	return &PatientGetRepository{db: db}
}

func (r *PatientGetRepository) GetPatientBasicInfo(userID uuid.UUID) (*databases.Patient, error) {
	var patient databases.Patient
	
	err := r.db.
		First(&patient, "user_id = ?", userID).Error

	if err != nil {
		return nil, err
	}

	return &patient, nil
}