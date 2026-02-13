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
	patientUsecase := _patientUsecases.NewPatientUsecase(tx, passwordSvc)
	patientReadRepo := _patientRepositories.NewPatientReadRepository(s.Db)
	patientGetUsecase := _patientUsecases.NewPatientGetUsecase(patientReadRepo)
	_patientControllers.NewPatientController(patientsGroup, patientUsecase, patientGetUsecase)

	diseaseGroupUsecase := _patientUsecases.NewDiseaseUsecase(tx)
	_patientControllers.NewDiseaseController(diseaseGroup, diseaseGroupUsecase)

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

