package repositories

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) entities.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateWithUser(req *entities.AdminCreateReq, hashedPass string) (*entities.AdminCreateRes, error) {
	returnValue := &entities.AdminCreateRes{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		user := databases.User{
			Username: req.Username,
			Password: hashedPass,
			Role:     "admin",
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		admin := databases.Admin{
			Title:          req.Title,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			UserID:         user.ID,
			CreatedBy:      req.CreatedBy,
			UpdatedBy:      req.CreatedBy,
		}

		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("create admin: %w", err)
		}

		returnValue = &entities.AdminCreateRes{
			ID:       admin.ID,
			FullName: fmt.Sprintf("%s%s %s", admin.Title, admin.FirstName, admin.LastName),
			Username: user.Username,
			Role:     user.Role,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return returnValue, nil
}