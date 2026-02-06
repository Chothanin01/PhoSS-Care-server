package repositories

import (

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type DiseaseRepository struct {
	db *gorm.DB
}

func NewDiseaseRepository(db *gorm.DB) *DiseaseRepository {
	return &DiseaseRepository{db: db}
}

func (r *DiseaseRepository) LinkPatientDiseases(patientID uuid.UUID, diseases []entities.PatientDiseaseEntity) error {
	var records []databases.PatientDisease
	for _, d := range diseases {
		records = append(records, databases.PatientDisease{
			PatientID: patientID,
			DiseaseID: d.DiseaseID,
			Disease: databases.Disease{
				Name: d.Name,
			},
		})
	}
	return r.db.Create(&records).Error
}
