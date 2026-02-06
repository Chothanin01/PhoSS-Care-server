package repositories

import (
	"fmt"
	"strconv"
	"time"
	"errors"

	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
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

type PatientReadRepository struct {
	db *gorm.DB
}

func NewTransactionGorm(db *gorm.DB) *TransactionGorm {
	return &TransactionGorm{db: db}
}

func (t *TransactionGorm) Do(fn func(entities.RepositorySet) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		repos := entities.RepositorySet{
			UserRepo:     NewUserRepository(tx),
			PatientRepo:  NewPatientRepository(tx),
			RelativeRepo: NewRelativeRepository(tx),
			DiseaseRepo:  NewDiseaseRepository(tx),
		}
		return fn(repos)
	})
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func NewPatientReadRepository(db *gorm.DB) *PatientReadRepository {
	return &PatientReadRepository{db: db}
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *PatientRepository) GenerateNextHnID() (string, error) {
	var last databases.Patient
	err := r.db.Order("hn_id DESC").First(&last).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "0000001", nil
	}
	lastInt, _ := strconv.Atoi(last.HnID)
	return fmt.Sprintf("%07d", lastInt+1), nil
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
		DOB:         dob,
		HnID:        req.HnID,
		IDCard:      req.IDCard,
		Rights:      req.Rights,
		Nationality: req.Nationality,
		Ethnicity:   req.Ethnicity,
		PhoneNumber: req.PhoneNumber,
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

func (r *PatientReadRepository) GetPatients(page, limit int) ([]databases.Patient, error) {
	var patients []databases.Patient
	offset := (page - 1) * limit
	err := r.db.
		Preload("Diseases.Disease").
		Preload("Appointments", "status = ?", "ongoing").
		Offset(offset).
		Limit(limit).
		Find(&patients).Error
	return patients, err
}

func (r *PatientReadRepository) GetPatientsWithFilter(req entities.PatientQueryParams) ([]databases.Patient, error) {
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

func (r *PatientReadRepository) CountPatients() (int64, error) {
	var count int64
	err := r.db.Model(&databases.Patient{}).Count(&count).Error
	return count, err
}

func (r *PatientReadRepository) CountPatientsWithFilter(req entities.PatientQueryParams) (int64, error) {
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
