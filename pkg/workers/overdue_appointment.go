package workers

import (
	"fmt"
	"time"

	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

func RunOverdueAppointmentWorker(db *gorm.DB) {
	fmt.Println("[CRON] Running Overdue Appointment Check...")

	today := time.Now().Format("2006-01-02")

	result := db.Model(&databases.Appoint{}).
		Where("DATE(date) = ? AND status IN ?", today, []string{"ongoing"}).
		Updates(map[string]interface{}{
			"status": "overdue", 
		})

	if result.Error != nil {
		fmt.Println("[CRON] Error updating overdue appointments:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		fmt.Println("[CRON] No overdue appointments found for", today)
		return
	}

	fmt.Printf("[CRON] Success! Marked %d appointments as overdue for %s\n", result.RowsAffected, today)
}