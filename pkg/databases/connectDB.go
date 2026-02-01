package databases

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/configs"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func SetupDatabaseConnection(cfg *configs.Config) (*gorm.DB, error) {

	dsn, err := utils.ConnectionUrlBuilder("gorm", cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to build database DSN: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect database: %w", err)
	}

	fmt.Println("Database connection established successfully")
	return db, nil
}
