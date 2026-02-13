package main

import (
	"log"
	"time"
	"fmt"

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
	var now time.Time
	db.Raw("SELECT now()").Scan(&now)
	fmt.Println("DB Time:", now)


	databases.MigrateAllIfEmpty(db)

	server := servers.NewServer(appconfig, db)
	server.Start()
}


