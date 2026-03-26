package repositories

import (
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"	
)

type GetAppointmentRepository struct {
	db *gorm.DB
}

func NewGetAppointmentRepository(db *gorm.DB) *GetAppointmentRepository {
	return &GetAppointmentRepository{db: db}
}

func (r *GetAppointmentRepository) GetAppointmentDetail(patientID uuid.UUID, diseaseID uuid.UUID) (*databases.Appoint, *databases.Admin, error) {
	var appoint databases.Appoint
	var admin databases.Admin

	err := r.db.
		Preload("Disease").
        Where("patient_id = ? AND disease_id = ? AND status IN ?", patientID, diseaseID, []string{"ongoing", "delay"}).
        First(&appoint).Error
    if err != nil {
        return nil, nil, err
    }
	
	if appoint.CreatedBy != nil {
        err = r.db.Where("user_id = ?", *appoint.CreatedBy).First(&admin).Error
        if err != nil {
            return &appoint, nil, nil 
        }
    }

	return &appoint, &admin, nil
}