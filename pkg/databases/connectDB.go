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

	fmt.Println("Database connection established successfully (Timezone: Asia/Bangkok)")
	return db, nil
}

func MigrateAllIfEmpty(db *gorm.DB) {
	var tableCount int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount).Error
	if err != nil {
		log.Fatalf("Failed to check existing tables: %v", err)
	}

	if tableCount > 0 {
		fmt.Println("Tables already exist — skipping migration.")
		return
	}

	fmt.Println("No tables found. Running initial migrations...")

	if err := db.AutoMigrate(
		&User{},
		&Admin{},
		&Patient{},
		&Health{},
		&Relative{},
		&Appoint{},
		&Disease{},
		&PatientDisease{},
		&Request{},
		&Vaccine{},
		&VaccinationRecord{},
		&Notification{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("All tables created successfully.")
}
