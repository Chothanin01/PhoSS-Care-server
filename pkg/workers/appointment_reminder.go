// 📂 workers/appointment_reminder.go
package workers

import (
	"fmt"
	"time"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

func RunAppointmentReminderWorker(db *gorm.DB) {
	fmt.Println("[CRON] Running Appointment Reminder Check...")

	targetDate := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	var upcomingAppoints []databases.Appoint

	err := db.Preload("CreatedByUser.Admin").
		Where("DATE(date) = ? AND status = ?", targetDate, "ongoing").
		Find(&upcomingAppoints).Error

	if err != nil || len(upcomingAppoints) == 0 {
		fmt.Println("[CRON] No upcoming appointments found for", targetDate)
		return
	}

	var notifications []databases.Notification

	for _, app := range upcomingAppoints {
		
		var existingCount int64
		db.Model(&databases.Notification{}).
			Where("appoint_id = ? AND header = ?", app.ID, "คุณมีนัดในอีก 2 วันข้างหน้า").
			Count(&existingCount)

		if existingCount > 0 {
			continue
		}

		doctorName := app.Doctor
		if app.CreatedByUser != nil && app.CreatedByUser.Admin != nil {
			doctorName = app.CreatedByUser.Admin.Title + app.CreatedByUser.Admin.FirstName 
		}

		bodyText := fmt.Sprintf("คุณมีนัดหมายกับ%s ในอีก 2 วันในเวลา %s", 
			doctorName, app.StartTime)

		notifications = append(notifications, databases.Notification{
			Header:    "คุณมีนัดในอีก 2 วันข้างหน้า",
			Body:      bodyText,
			PatientID: app.PatientID,
			AppointID: &app.ID, 
		})
	}

	if len(notifications) == 0 {
		fmt.Println("[CRON] All upcoming appointments already have reminders. Nothing new to save.")
		return
	}

	if err := db.Create(&notifications).Error; err != nil {
		fmt.Println("[CRON] Error saving notifications:", err)
		return
	}

	fmt.Printf("[CRON] Success! Created %d reminders for %s\n", len(notifications), targetDate)
}