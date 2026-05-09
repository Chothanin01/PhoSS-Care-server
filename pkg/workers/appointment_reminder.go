package workers

import (
    "fmt"
    "time"

    "github.com/chothanin01/PhoSS-Care-server/pkg/databases"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

func RunAppointmentReminderWorker(db *gorm.DB) {

    loc, _ := time.LoadLocation("Asia/Bangkok")
    targetDate := time.Now().In(loc).AddDate(0, 0, 2).Format("2006-01-02")

    fmt.Printf("[CRON] Running Appointment Reminder Check for %s...\n", targetDate)

    var upcomingAppoints []databases.Appoint

    err := db.Preload("Doctor").
        Preload("CreatedByUser.Admin").
        Where("DATE(date) = ? AND status = ?", targetDate, "ongoing").
        Find(&upcomingAppoints).Error

    if err != nil || len(upcomingAppoints) == 0 {
        fmt.Println("[CRON] No upcoming appointments found for", targetDate)
        return
    }

    var appointIDs []uuid.UUID
    for _, app := range upcomingAppoints {
        appointIDs = append(appointIDs, app.ID)
    }

    var existingNotifs []databases.Notification
    db.Where("appoint_id IN ? AND header = ?", appointIDs, "คุณมีนัดในอีก 2 วันข้างหน้า").
        Find(&existingNotifs)

    existingMap := make(map[uuid.UUID]bool)
    for _, notif := range existingNotifs {
        if notif.AppointID != nil {
            existingMap[*notif.AppointID] = true
        }
    }

    var notifications []databases.Notification

    for _, app := range upcomingAppoints {
        
        if existingMap[app.ID] {
            continue 
        }

        doctorName := app.Doctor.Title + app.Doctor.FirstName + " " + app.Doctor.LastName 
        if app.CreatedByUser != nil && app.CreatedByUser.Admin != nil {
            doctorName = app.CreatedByUser.Admin.Title + app.CreatedByUser.Admin.FirstName 
        }

        bodyText := fmt.Sprintf("คุณมีนัดหมายกับ %s ในอีก 2 วันในเวลา %s", doctorName, app.StartTime)

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