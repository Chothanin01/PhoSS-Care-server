package entities

type Disease struct {
	DiseaseID uint   `json:"disease_id"`
	Name      string `json:"name"`
}

type DiseaseWithStatus struct {
	DiseaseID uint   `json:"disease_id"`
	Name           string `json:"name"`
	HasAppointment bool   `json:"has_appointment"`
}

type DiseaseRepository interface {
	LinkPatientDiseases(patientID uint, diseases []PatientDiseaseEntity) error
}

