package repositories

import (
	"fmt"
	"time"
	"strconv"

	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
	"github.com/chothanin01/PhoSS-Care-server/pkg/databases"
	"gorm.io/gorm"
)

type TransactionGorm struct {
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

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) GenerateNextHnID() (string, error) {
	var last databases.Patient
	if err := r.db.Order("hn_id desc").First(&last).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "0000001", nil
		}
		return "", err
	}

	lastInt, _ := strconv.Atoi(last.HnID)
	return fmt.Sprintf("%07d", lastInt+1), nil
}


func (r *PatientRepository) CreateWithUser(req *entities.PatientCreateReq, userID uint) (*entities.PatientCreateRes, error) {
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
			HouseNumber: req.Address.HouseNumber,
			VillageNumber:         req.Address.VillageNumber,
			Alley:         req.Address.Alley,
			Road:        req.Address.Road,
			SubDistrict: req.Address.SubDistrict,
			District:    req.Address.District,
			Province:    req.Address.Province,
			ZipCode:     req.Address.ZipCode,
		},
		Allergy: req.Allergy,
		UserID:  userID,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.UpdatedBy,
	}

	if err := r.db.Create(&patient).Error; err != nil {
		return nil, err
	}

	return &entities.PatientCreateRes{
		Id:          uint64(patient.ID),
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

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
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

type RelativeRepository struct {
	db *gorm.DB
}

func NewRelativeRepository(db *gorm.DB) *RelativeRepository {
	return &RelativeRepository{db: db}
}

func (r *RelativeRepository) Create(relatives []entities.RelativeEntity) error {
	var records []databases.Relative
	for _, rel := range relatives {
		records = append(records, databases.Relative{
			Title:       rel.Title,
			FirstName:   rel.FirstName,
			LastName:    rel.LastName,
			PhoneNumber: rel.PhoneNumber,
			Role:        rel.Role,
			PatientID:   rel.PatientID,
			Address: databases.Address{
				HouseNumber:  rel.Address.HouseNumber,
				VillageNumber: rel.Address.VillageNumber,
				Alley:         rel.Address.Alley,
				Road:         rel.Address.Road,
				SubDistrict:  rel.Address.SubDistrict,
				District:     rel.Address.District,
				Province:     rel.Address.Province,
				ZipCode:      rel.Address.ZipCode,
			},
			CreatedBy: rel.CreatedBy,
			UpdatedBy: rel.UpdatedBy,
		})
	}
	return r.db.Create(&records).Error
}

type DiseaseRepository struct {
	db *gorm.DB
}

func NewDiseaseRepository(db *gorm.DB) *DiseaseRepository {
	return &DiseaseRepository{db: db}
}

func (r *DiseaseRepository) LinkPatientDiseases(patientID uint, diseases []entities.PatientDiseaseEntity) error {
	var records []databases.PatientDisease
	for _, d := range diseases {
		records = append(records, databases.PatientDisease{
			PatientID: patientID,
			DiseaseID: d.DiseaseID,
			Disease: databases.Disease{
				Name: d.Name,
			},
		})
	}
	return r.db.Create(&records).Error
}
