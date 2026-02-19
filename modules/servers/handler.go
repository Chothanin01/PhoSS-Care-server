package servers

import (
	_patientControllers "github.com/chothanin01/PhoSS-Care-server/modules/patients/controllers"
	_patientRepositories "github.com/chothanin01/PhoSS-Care-server/modules/patients/repositories"
	_patientUsecases "github.com/chothanin01/PhoSS-Care-server/modules/patients/usecases"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")
	patientsGroup := v1.Group("/patients")
	diseaseGroup := v1.Group("/diseases")

	passwordSvc := utils.NewPasswordService()
	tx := _patientRepositories.NewTransactionGorm(s.Db)

	// Usecases
	patientUsecase := _patientUsecases.NewPatientUsecase(tx, passwordSvc)
	patientGetRepo := _patientRepositories.NewPatientGetRepository(s.Db)
	patientGetUsecase := _patientUsecases.NewPatientGetUsecase(patientGetRepo)

	// 🔹 New relative usecase
	relativeRepo := _patientRepositories.NewRelativeRepository(s.Db)
	relativeUsecase := _patientUsecases.NewRelativeUsecase(relativeRepo)

	// Controllers
	_patientControllers.NewPatientController(patientsGroup, patientUsecase, patientGetUsecase, patientUsecase)
	_patientControllers.NewRelativeController(patientsGroup, relativeUsecase)

	diseaseGetRepo := _patientRepositories.NewDiseaseRepository(s.Db)
	diseaseGroupUsecase := _patientUsecases.NewDiseaseUsecase(diseaseGetRepo)
	_patientControllers.NewDiseaseController(diseaseGroup, diseaseGroupUsecase)



	// Default 404 handler
	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusNotFound,
			"message":     "endpoint not found",
			"result":      nil,
		})
	})

	return nil
}


