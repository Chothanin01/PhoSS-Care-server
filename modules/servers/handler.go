package servers

import (
	_patientHandlers "github.com/chothanin01/PhoSS-Care-server/modules/patients/controllers"
	_patientRepositories "github.com/chothanin01/PhoSS-Care-server/modules/patients/repositories"
	_patientUsecases "github.com/chothanin01/PhoSS-Care-server/modules/patients/usecases"

	"github.com/gofiber/fiber/v2"
)

// MapHandlers maps all HTTP handlers to routes
func (s *Server) MapHandlers() error {
	// API version group
	v1 := s.App.Group("/v1")

	// ===== PATIENTS =====
	patientsGroup := v1.Group("/patients")
	patientRepo := _patientRepositories.NewPatientRepository(s.Db)
	patientUsecase := _patientUsecases.NewPatientUsecase(patientRepo)
	_patientHandlers.NewPatientController(patientsGroup, patientUsecase)

	// ===== FALLBACK (404 Handler) =====
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
