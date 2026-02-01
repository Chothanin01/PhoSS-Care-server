package main

import (
	"fmt"
	"log"

	"github.com/chothanin01/PhoSS-Care-server/configs"
	"github.com/chothanin01/PhoSS-Care-server/modules/servers"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appconfig := configs.LoadConfig()

	db, err := databases.SetupDatabaseConnection(appconfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connection successful:", db)

	db.AutoMigrate(&databases.Admin{}, &databases.Patient{}, &databases.User{}, &databases.Health{}, &databases.Relative{},
		&databases.Appoint{}, &databases.Disease{}, &databases.PatientDisease{}, &databases.Notification{}, 
		&databases.Vaccine{}, &databases.VaccinationRecord{}, &databases.Request{})

	server := servers.NewServer(appconfig, db)
	server.Start()
}


