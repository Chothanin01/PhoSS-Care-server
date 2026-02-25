package repositories

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type RelativeRepository struct {
	db *gorm.DB
}

func NewRelativeRepository(db *gorm.DB) *RelativeRepository {
	return &RelativeRepository{db: db}
}

func (r *RelativeRepository) Create(relatives []entities.RelativeEntity) error {
	var records []databases.Relative
	for _, rel := range relatives {
		records = append(records, databases.Relative{
			Title:       rel.Title,
			FirstName:   rel.FirstName,
			LastName:    rel.LastName,
			PhoneNumber: rel.PhoneNumber,
			Role:        rel.Role,
			PatientID:   rel.PatientID,
			Address: databases.Address{
				HouseNumber:  rel.Address.HouseNumber,
				VillageNumber: rel.Address.VillageNumber,
				Alley:         rel.Address.Alley,
				Road:         rel.Address.Road,
				SubDistrict:  rel.Address.SubDistrict,
				District:     rel.Address.District,
				Province:     rel.Address.Province,
				ZipCode:      rel.Address.ZipCode,
			},
			CreatedBy: rel.CreatedBy,
			UpdatedBy: rel.UpdatedBy,
		})
	}
	return r.db.Create(&records).Error
}

