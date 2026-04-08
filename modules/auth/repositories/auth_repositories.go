package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/auth/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) entities.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetByUsername(username string) (*entities.User, error) {
	var user databases.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &entities.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}, nil
}

func (r *authRepository) GetPatientDiseases(userID uuid.UUID) ([]map[string]interface{}, error) {
	type DiseaseInfo struct {
		DiseaseID   uuid.UUID
		DiseaseName string
	}

	var diseases []DiseaseInfo
	err := r.db.Table("patient_disease pd").
		Select("d.id as disease_id, d.name as disease_name").
		Joins("JOIN disease d ON d.id = pd.disease_id").
		Joins("JOIN patient p ON p.id = pd.patient_id").
		Where("p.user_id = ?", userID).
		Scan(&diseases).Error

	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(diseases))
	for i, d := range diseases {
		result[i] = map[string]interface{}{
			"disease_id":   d.DiseaseID,
			"disease_name": d.DiseaseName,
		}
	}

	return result, nil
}