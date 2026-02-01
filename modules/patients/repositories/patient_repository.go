package repositories

import (
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
	"time"
	"fmt"
)

type PatientRepository struct {
	Db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{
		Db: db,
	}
}

func (r *PatientRepository) Create(req *entities.PatientCreateReq) (*entities.PatientCreateRes, error) {
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return nil, fmt.Errorf("invalid dob format: %v", err)
	}
	
	patient := databases.Patient{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DOB:            dob,
		HnID:           req.HnID,
		IDCard:         req.IDCard,
		Rights:         req	.Rights,
	}
	if err := r.Db.Create(&patient).Error; err != nil {
		return nil, err
	}

	res := &entities.PatientCreateRes{
		Id:          uint64(patient.ID),
		FirstName:   patient.FirstName,
		LastName:    patient.LastName,
		HnID:        patient.HnID,
		IDCard:      patient.IDCard,
		PhoneNumber: patient.PhoneNumber,
		Rights:      patient.Rights,
	}

	return res, nil
}
