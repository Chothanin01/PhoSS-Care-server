package workers

import (
	"fmt"
	"log"


	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)


func StartCronJobs(db *gorm.DB) {
	c := cron.New()

	_, err := c.AddFunc("0 8 * * *", func() {
		RunAppointmentReminderWorker(db)
	})
	if err != nil {
		log.Fatalf("Failed to setup appointment reminder job: %v", err)
	}

	_, err = c.AddFunc("0 15 * * *", func() {
		RunOverdueAppointmentWorker(db)
	})
	if err != nil {
		log.Fatalf("Failed to setup overdue appointment job: %v", err)
	}

	c.Start()
	fmt.Println("Background Cron scheduler started successfully.")
}