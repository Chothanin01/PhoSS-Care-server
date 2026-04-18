package repositories

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type vaccineQueryRepo struct {
	db *gorm.DB
}

func NewVaccineQueryRepo(db *gorm.DB) entities.VaccineQueryRepo {
	return &vaccineQueryRepo{db: db}
}

func (r *vaccineQueryRepo) buildBaseQuery(patientID uuid.UUID, filter string) *gorm.DB {
	
	patientRecords := r.db.Table("vaccination_record").
		Select("vaccination_record.vaccine_id, vaccination_record.status, appoint.date").
		Joins("JOIN appoint ON appoint.id = vaccination_record.appoint_id").
		Where("appoint.patient_id = ?", patientID)

	query := r.db.Table("vaccine").
		Select("vaccine.id, vaccine.name, vaccine.type, vaccine.age, pr.status as v_status, pr.date as v_date").
		Joins("LEFT JOIN (?) pr ON pr.vaccine_id = vaccine.id", patientRecords)

	switch filter {
	case "completed":
		query = query.Where("pr.status = ?", "completed") 
	case "ongoing":
		query = query.Where("pr.status = ?", "ongoing")
	case "not_vaccinated":
		query = query.Where("pr.status IS NULL")
	}

	return query
}

func (r *vaccineQueryRepo) CountVaccines(patientID uuid.UUID, filter string) (int64, error) {
	var count int64
	err := r.buildBaseQuery(patientID, filter).Count(&count).Error
	return count, err
}

func (r *vaccineQueryRepo) GetVaccineList(patientID uuid.UUID, filter string, limit int, offset int) ([]entities.VaccineItemEntity, error) {

	type Result struct {
		ID      uuid.UUID
		Name    string
		Type    string
		Age     string
		VStatus *string
		VDate   *time.Time
	}

	var results []Result

	err := r.buildBaseQuery(patientID, filter).
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	items := make([]entities.VaccineItemEntity, len(results))
	for i, res := range results {
		
		status := "not_vaccinated"
		if res.VStatus != nil {
			status = *res.VStatus
		}

		dateStr := ""
		if res.VDate != nil {
			dateStr = res.VDate.Format("2006-01-02")
		}

		items[i] = entities.VaccineItemEntity{
			VaccineID:        res.ID,
			Name:             res.Name,
			Type:             res.Type,
			Age:              res.Age,
			VaccinatedDate:   dateStr,
			VaccinatedStatus: status,
		}
	}

	return items, nil
}

func (r *vaccineQueryRepo) GetVaccineDetail(patientID uuid.UUID, vaccineID uuid.UUID) (*entities.VaccineDetailEntity, error) {
	// Temporary struct to catch the hybrid SQL result
	type Result struct {
		ID         uuid.UUID
		Name       string
		Type       string
		Age        string
		Effect     string
		Note       string
		VStatus    *string
		VDate      *time.Time
		VDose      *int
		VDoctor    *string
	}

	var res Result

	patientRecords := r.db.Table("vaccination_record").
		Select("vaccination_record.vaccine_id, vaccination_record.status, vaccination_record.dose_number, vaccination_record.vaccine_doctor, appoint.date").
		Joins("JOIN appoint ON appoint.id = vaccination_record.appoint_id").
		Where("appoint.patient_id = ?", patientID)

	err := r.db.Table("vaccine").
		Select("vaccine.id, vaccine.name, vaccine.type, vaccine.age, vaccine.effect, vaccine.note, pr.status as v_status, pr.date as v_date, pr.dose_number as v_dose, pr.vaccine_doctor as v_doctor").
		Joins("LEFT JOIN (?) pr ON pr.vaccine_id = vaccine.id", patientRecords).
		Where("vaccine.id = ?", vaccineID).
		First(&res).Error

	if err != nil {
		return nil, err 
	}
	status := "not_vaccinated"
	if res.VStatus != nil {
		status = *res.VStatus
	}

	dateStr := ""
	if res.VDate != nil {
		dateStr = res.VDate.Format("2006-01-02")
	}

	dose := 0
	if res.VDose != nil {
		dose = *res.VDose
	}

	doctor := ""
	if res.VDoctor != nil {
		doctor = *res.VDoctor
	}

	return &entities.VaccineDetailEntity{
		ID:               res.ID,
		Name:             res.Name,
		Type:             res.Type,
		Age:              res.Age,
		Effect:           res.Effect,
		Note:             res.Note,
		VaccinatedDate:   dateStr,
		VaccinatedStatus: status,
		DoseNumber:       dose,
		VaccineDoctor:    doctor,
	}, nil
}