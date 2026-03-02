package databases

import (
	"fmt"
	"os"

	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"gorm.io/gorm"
)

func SeedSuperAdmin(db *gorm.DB, passwordSvc utils.PasswordService) error {
	var count int64
	if err := db.Model(&User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing admins: %w", err)
	}

	if count > 0 {
		return nil
	}

	username := os.Getenv("SUPERADMIN_USERNAME")
	password := os.Getenv("SUPERADMIN_PASSWORD")
	firstname := os.Getenv("SUPERADMIN_FIRSTNAME")
	lastname := os.Getenv("SUPERADMIN_LASTNAME")

	if username == "" || password == "" {
		return fmt.Errorf("missing SUPERADMIN_USERNAME or SUPERADMIN_PASSWORD in .env")
	}

	hashed, err := passwordSvc.Hash(password)
	if err != nil {
		return fmt.Errorf("failed to hash superadmin password: %w", err)
	}

	user := User{
		Username: username,
		Password: hashed,
		Role:     "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create superadmin user: %w", err)
	}

	admin := Admin{
		Title:     "Mr.",
		FirstName: firstname,
		LastName:  lastname,
		UserID:    user.ID,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create superadmin profile: %w", err)
	}

	fmt.Printf("Superadmin created: %s (env-based)\n", username)
	return nil
}

func SeedDiseases(db *gorm.DB) error {
	var superadmin User
	superadminUsername := os.Getenv("SUPERADMIN_USERNAME")

	if superadminUsername == "" {
		return fmt.Errorf("missing SUPERADMIN_USERNAME env; cannot assign CreatedBy/UpdatedBy for diseases")
	}

	if err := db.Where("username = ?", superadminUsername).First(&superadmin).Error; err != nil {
		return fmt.Errorf("failed to find superadmin (username=%s): %w", superadminUsername, err)
	}

	defaultDiseases := []Disease{
		{Name: "โรคความดันโลหิตสูง"},
		{Name: "โรคเบาหวาน"},
		{Name: "วัณโรค"},
		{Name: "วัคซีน"},
	}

	for _, d := range defaultDiseases {
		var count int64
		if err := db.Model(&Disease{}).Where("name = ?", d.Name).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check disease %s: %w", d.Name, err)
		}

		if count == 0 {
			d.CreatedBy = &superadmin.ID
			d.UpdatedBy = &superadmin.ID
			if err := db.Create(&d).Error; err != nil {
				return fmt.Errorf("failed to seed disease %s: %w", d.Name, err)
			}
			fmt.Printf("Seeded disease: %s (by superadmin: %s)\n", d.Name, superadmin.Username)
		}
	}

	return nil
}