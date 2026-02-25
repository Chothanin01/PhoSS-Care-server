package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionGorm struct {
	db *gorm.DB
}

type UserRepository struct {
	db *gorm.DB
}

type PatientRepository struct {
	db *gorm.DB
}

type PatientGetRepository struct {
	db *gorm.DB
}

func NewTransactionGorm(db *gorm.DB) *TransactionGorm {
	return &TransactionGorm{db: db}
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func NewPatientGetRepository(db *gorm.DB) *PatientGetRepository {
	return &PatientGetRepository{db: db}
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (t *TransactionGorm) Do(fn func(entities.RepositorySet) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		repos := entities.RepositorySet{
			UserRepo:     NewUserRepository(tx),
			PatientRepo:  NewPatientRepository(tx),
			RelativeRepo: NewRelativeRepository(tx),
			DiseaseRepo:  NewDiseaseRepository(tx),
			PateintUpdateRepo: NewPatientRepository(tx),
		}
		return fn(repos)
	})
}

// ---------------------- CREATE ----------------------

func (r *PatientRepository) GenerateNextHnID() (string, error) {
	var lastHn string

	err := r.db.
		Model(&databases.Patient{}).
		Select("hn_id").
		Order("hn_id DESC").
		Limit(1).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Scan(&lastHn).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", fmt.Errorf("fetch last HN ID: %w", err)
	}

	if lastHn == "" {
		return "0000001", nil
	}

	lastHn = strings.TrimSpace(lastHn)

	var num int
	_, err = fmt.Sscanf(lastHn, "%d", &num)
	if err != nil {
		return "", fmt.Errorf("invalid HN ID format (%s): %w", lastHn, err)
	}

	nextHn := fmt.Sprintf("%07d", num+1)
	return nextHn, nil
}

func (r *PatientRepository) CreateWithUser(req *entities.PatientCreateReq, userID uuid.UUID) (*entities.PatientCreateRes, error) {
	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return nil, fmt.Errorf("invalid dob format: %v", err)
	}


	patient := databases.Patient{
		Title:       req.Title,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Sex: 	     req.Sex,
		DOB:         dob,
		HnID:        req.HnID,
		IDCard:      req.IDCard,
		Rights:      req.Rights,
		Nationality: req.Nationality,
		Ethnicity:   req.Ethnicity,
		PhoneNumber: req.PhoneNumber,
		Weight: 	 req.Weight,
		Height: 	 req.Height,
		Address: databases.Address{
			HouseNumber:   req.Address.HouseNumber,
			VillageNumber: req.Address.VillageNumber,
			Alley:         req.Address.Alley,
			Road:          req.Address.Road,
			SubDistrict:   req.Address.SubDistrict,
			District:      req.Address.District,
			Province:      req.Address.Province,
			ZipCode:       req.Address.ZipCode,
		},
		Allergy:   req.Allergy,
		UserID:    userID,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
	}

	if err := r.db.Create(&patient).Error; err != nil {
		return nil, err
	}
	fmt.Println("Creating patient at:", time.Now())
	fmt.Printf("Patient.CreatedAt before insert: %v\n", patient.CreatedAt)


	for _, d := range req.Diseases {
	var disease databases.Disease
	if err := r.db.First(&disease, "id = ?", d.DiseaseID).Error; err != nil {
		return nil, fmt.Errorf("disease not found: %v", d.DiseaseID)
	}

	link := databases.PatientDisease{
		PatientID: patient.ID,
		DiseaseID: d.DiseaseID,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
	}
	if err := r.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(&link).Error; err != nil {
		return nil, err
		}
	}


	return &entities.PatientCreateRes{
		Id:          patient.ID,
		FirstName:   patient.FirstName,
		LastName:    patient.LastName,
		HnID:        patient.HnID,
		IDCard:      patient.IDCard,
		PhoneNumber: patient.PhoneNumber,
		Rights:      patient.Rights,
		Nationality: patient.Nationality,
		Ethnicity:   patient.Ethnicity,
	}, nil
}


func (r *UserRepository) Create(username, password, role string) (*databases.User, error) {
	user := &databases.User{
		Username: username,
		Password: password,
		Role:     role,
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// ---------------------- GET ----------------------

func (r *PatientGetRepository) GetPatients(page, limit int) ([]databases.Patient, error) {
	var patients []databases.Patient
	offset := (page - 1) * limit
	err := r.db.
		Preload("Diseases.Disease").
		Preload("Appointments", "status = ?", "Ongoing").
		Offset(offset).
		Limit(limit).
		Find(&patients).Error
	return patients, err
}

func (r *PatientGetRepository) GetPatientsWithFilter(req entities.PatientQueryParams) ([]databases.Patient, error) {
	var patients []databases.Patient
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	query := r.db.Model(&databases.Patient{}).
		Preload("Diseases.Disease").
		Preload("Appointments")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where(
			"hn_id ILIKE ? OR id_card ILIKE ? OR CONCAT(title, ' ', first_name, ' ', last_name) ILIKE ?",
			search, search, search,
		)
	}

	if len(req.Diseases) > 0 {
		query = query.Joins("JOIN patient_disease pd ON pd.patient_id = patient.id").
			Joins("JOIN disease d ON d.id = pd.disease_id").
			Where("d.name IN ?", req.Diseases)
	}

	if req.Appoint != nil {
		if *req.Appoint {
			query = query.Joins("JOIN appoint a ON a.patient_id = patient.id AND a.status = ?", "ongoing")
		} else {
			query = query.Where("NOT EXISTS (SELECT 1 FROM appoint a WHERE a.patient_id = patient.id AND a.status = ?)", "ongoing")
		}
	}

	err := query.Offset(offset).Limit(req.Limit).Find(&patients).Error
	return patients, err
}

func (r *PatientGetRepository) CountPatients() (int64, error) {
	var count int64
	err := r.db.Model(&databases.Patient{}).Count(&count).Error
	return count, err
}

func (r *PatientGetRepository) CountPatientsWithFilter(req entities.PatientQueryParams) (int64, error) {
	var count int64

	query := r.db.Model(&databases.Patient{})

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where(
			"hn_id ILIKE ? OR id_card ILIKE ? OR CONCAT(title, ' ', first_name, ' ', last_name) ILIKE ?",
			search, search, search,
		)
	}

	if len(req.Diseases) > 0 {
		query = query.Joins("JOIN patient_disease pd ON pd.patient_id = patient.id").
			Joins("JOIN disease d ON d.id = pd.disease_id").
			Where("d.name IN ?", req.Diseases)
	}

	if req.Appoint != nil {
		if *req.Appoint {
			query = query.Joins("JOIN appoint a ON a.patient_id = patient.id AND a.status = ?", "ongoing")
		} else {
			query = query.Where("NOT EXISTS (SELECT 1 FROM appoint a WHERE a.patient_id = patient.id AND a.status = ?)", "ongoing")
		}
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *PatientGetRepository) GetPatientInfoByID(id uuid.UUID) (*databases.Patient, error) {
	var patient databases.Patient
	err := r.db.
		Preload("Diseases.Disease").
		Preload("Relatives").
		First(&patient, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *PatientGetRepository) GetPatientDiseasesInfoByID(patientID, diseaseID uuid.UUID) (*databases.Patient, error) {
	var patient databases.Patient

	err := r.db.
		Preload("Diseases.Disease", "id = ?", diseaseID).
		Preload("Appointments", func(db *gorm.DB) *gorm.DB {
			return db.Where("disease_id = ?", diseaseID).Order("no DESC")
		}).
		Preload("Healths").
		First(&patient, "id = ?", patientID).Error

	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *PatientGetRepository) GetPatientAppointmentsByID(patientID uuid.UUID) (*databases.Patient, error) {
	var patient databases.Patient

	err := r.db.
		Preload("Appointments", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("status IN ?", []string{"Ongoing", "Delay"}).
				Order("no DESC")
		}).
		Preload("Appointments.Disease").
		First(&patient, "id = ?", patientID).Error

	if err != nil {
		return nil, err
	}

	return &patient, nil
}

// ---------------------- EDIT ----------------------

func (r *PatientRepository) UpdatePatientInfo(id uuid.UUID, req *entities.PatientUpdateReq) (*entities.PatientUpdateRes, error) {
	var patient databases.Patient

	if err := r.db.First(&patient, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("patient not found: %w", err)
	}

	var dob time.Time
	if req.DOB != "" {
		parsedDOB, err := time.Parse("2006-01-02", req.DOB)
		if err != nil {
			return nil, fmt.Errorf("invalid dob format: %v", err)
		}
		dob = parsedDOB
	}

	if req.IDCard != "" && req.IDCard != patient.IDCard {
		var exists bool
		r.db.Model(&databases.User{}).
			Where("username = ? AND id <> ?", req.IDCard, patient.UserID).
			Select("count(*) > 0").
			Find(&exists)
		if exists {
			return nil, fmt.Errorf("ID card already used by another user")
		}

		if err := r.db.Model(&databases.User{}).
			Where("id = ?", patient.UserID).
			Updates(map[string]interface{}{
				"username":   req.IDCard,
				"updated_by": req.UpdatedBy,   
				"updated_at": time.Now(),    
			}).Error; err != nil {
			return nil, fmt.Errorf("update user username: %w", err)
		}
	}

	patient.Title = req.Title
	patient.FirstName = req.FirstName
	patient.LastName = req.LastName
	patient.Sex = req.Sex
	if !dob.IsZero() {
		patient.DOB = dob
	}
	patient.Weight = req.Weight
	patient.Height = req.Height
	patient.IDCard = req.IDCard
	patient.Rights = req.Rights
	patient.Nationality = req.Nationality
	patient.Ethnicity = req.Ethnicity
	patient.PhoneNumber = req.PhoneNumber
	patient.Address = databases.Address{
		HouseNumber:   req.Address.HouseNumber,
		VillageNumber: req.Address.VillageNumber,
		Alley:         req.Address.Alley,
		Road:          req.Address.Road,
		SubDistrict:   req.Address.SubDistrict,
		District:      req.Address.District,
		Province:      req.Address.Province,
		ZipCode:       req.Address.ZipCode,
	}
	patient.UpdatedBy = req.UpdatedBy

	if err := r.db.Save(&patient).Error; err != nil {
		return nil, fmt.Errorf("update patient info: %w", err)
	}

	res := &entities.PatientUpdateRes{
		ID:          patient.ID,
		Fullname:    fmt.Sprintf("%s%s %s", patient.Title, patient.FirstName, patient.LastName),
		HnNumber:    patient.HnID,
		IDCard:      patient.IDCard,
		PhoneNumber: patient.PhoneNumber,
		Rights:      patient.Rights,
		Nationality: patient.Nationality,
		Ethnicity:   patient.Ethnicity,
	}

	return res, nil
}
