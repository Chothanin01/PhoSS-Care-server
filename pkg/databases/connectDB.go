package databases

import (
	"fmt"
	"log"

	"github.com/chothanin01/PhoSS-Care-server/configs"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func SetupDatabaseConnection(cfg *configs.Config) (*gorm.DB, error) {
	dsn, err := utils.ConnectionUrlBuilder("gorm", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build database DSN: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := db.Exec("SET TIME ZONE 'Asia/Bangkok';").Error; err != nil {
		return nil, fmt.Errorf("failed to set timezone: %w", err)
	}

	fmt.Println("Database connection established successfully.")
	
	if err := MigrateAll(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	passwordSvc := utils.NewPasswordService()
	if err := SeedSuperAdmin(db, passwordSvc); err != nil {
		log.Fatalf("Failed to seed superadmin: %v", err)
	}

	if err := SeedDiseases(db); err != nil {
		log.Fatalf("Failed to seed diseases: %v", err)
	} 
	
	if err := SeedVaccines(db); err != nil {
		log.Fatalf("Failed to seed vaccines: %v", err)
	} 


	return db, nil
}

func MigrateAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.AutoMigrate(&User{}, &Admin{}); err != nil {
			return fmt.Errorf("migrate user/admin failed: %w", err)
		}

		if err := tx.AutoMigrate(&Patient{}, &Disease{}, &PatientDisease{}); err != nil {
			return fmt.Errorf("migrate patient/disease failed: %w", err)
		}

		if err := tx.AutoMigrate(&Appoint{}, &Relative{}); err != nil {
			return fmt.Errorf("migrate appointment-related failed: %w", err)
		}

		if err := tx.AutoMigrate(&Vaccine{}, &VaccinationRecord{}); err != nil {
			return fmt.Errorf("migrate vaccine failed: %w", err)
		}

		if err := tx.AutoMigrate(&Request{}, &Notification{}); err != nil {
			return fmt.Errorf("migrate request/notification failed: %w", err)
		}

		return nil
	})
}

