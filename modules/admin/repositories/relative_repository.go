package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type RelativeRepository struct {
	db *gorm.DB
}

func NewRelativeRepository(db *gorm.DB) *RelativeRepository {
	return &RelativeRepository{db: db}
}

func (r *RelativeRepository) Create(relatives []entities.RelativeEntity, adminID uuid.UUID) error {
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
			CreatedBy: &adminID,
			UpdatedBy: &adminID,
		})
	}
	return r.db.Create(&records).Error
}

func (r *RelativeRepository) UpdateRelativeInfo(patientID uuid.UUID, req *entities.RelativeAllUpdateReq, adminID uuid.UUID) (*entities.RelativeAllUpdateRes, error) {
	roles := map[string]entities.RelativeUpdateReq{
		"kin":       req.Kin,
		"caretaker": req.Caretaker,
		"medicine":  req.Medicine,
	}

	results := &entities.RelativeAllUpdateRes{}

	for role, data := range roles {
		if data.FirstName == "" && data.LastName == "" {
			continue
		}

		var relative databases.Relative
		if err := r.db.First(&relative, "patient_id = ? AND role = ?", patientID, role).Error; err != nil {
			return nil, fmt.Errorf("%s not found: %w", role, err)
		}

		update := map[string]interface{}{
			"title":        data.Title,
			"first_name":   data.FirstName,
			"last_name":    data.LastName,
			"phone_number": data.PhoneNumber,
			"address":      data.Address,
			"updated_by":   adminID,
			"updated_at":   time.Now(),
		}

		if err := r.db.Model(&relative).Updates(update).Error; err != nil {
			return nil, fmt.Errorf("update %s failed: %w", role, err)
		}

		res := entities.RelativeUpdateRes{
			ID:          relative.ID,
			FullName:    fmt.Sprintf("%s%s %s", data.Title, data.FirstName, data.LastName),
			PhoneNumber: data.PhoneNumber,
			Role:        role,
		}

		switch role {
		case "kin":
			results.Kin = res
		case "caretaker":
			results.Caretaker = res
		case "medicine":
			results.Medicine = res
		}
	}

	return results, nil
}

func (r *RelativeRepository) UpdateOfficerInfo(patientID uuid.UUID, req *entities.OfficerAllUpdateReq, adminID uuid.UUID) (*entities.OfficerAllUpdateRes, error) {
	var result entities.OfficerAllUpdateRes

	roles := map[string]entities.OfficerUpdateReq{
		"house": req.House,
		"nurse": req.Nurse,
	}

	for role, data := range roles {
		if data.FirstName == "" && data.LastName == "" {
			continue
		}

		fullname := fmt.Sprintf("%s%s %s", data.Title, data.FirstName, data.LastName)

		updateData := map[string]interface{}{
			"title":        data.Title,
			"first_name":   data.FirstName,
			"last_name":    data.LastName,
			"updated_by":   adminID,
		}

		if err := r.db.Model(&databases.Relative{}).
			Where("patient_id = ? AND role = ?", patientID, role).
			Updates(updateData).Error; err != nil {
			return nil, fmt.Errorf("update %s officer: %w", role, err)
		}

		switch role {
		case "house":
			result.House = entities.OfficerUpdateRes{	
				ID:       uuid.New(),
				FullName: fullname,
				Role:     "house",
			}
		case "nurse":
			result.Nurse = entities.OfficerUpdateRes{
				ID:       uuid.New(),
				FullName: fullname,
				Role:     "nurse",
			}
		}
	}

	return &result, nil
}