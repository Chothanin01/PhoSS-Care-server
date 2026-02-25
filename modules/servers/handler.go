package servers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	_adminControllers "github.com/chothanin01/PhoSS-Care-server/modules/admin/controllers"
	_adminRepositories "github.com/chothanin01/PhoSS-Care-server/modules/admin/repositories"
	_adminUsecases "github.com/chothanin01/PhoSS-Care-server/modules/admin/usecases"

	_authControllers "github.com/chothanin01/PhoSS-Care-server/modules/auth/controllers"
	_authRepositories "github.com/chothanin01/PhoSS-Care-server/modules/auth/repositories"
	_authUsecases "github.com/chothanin01/PhoSS-Care-server/modules/auth/usecases"

	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	// ------------------ 🔐 AUTH SETUP ------------------
	passSvc := utils.NewPasswordService()
	authRepo := _authRepositories.NewAuthRepository(s.Db)

	jwtAdmin := utils.NewJWTService(os.Getenv("JWT_SECRET_ADMIN"), "admin-api")
	jwtPatient := utils.NewJWTService(os.Getenv("JWT_SECRET_PATIENT"), "patient-api")

	adminAuthUC := _authUsecases.NewAuthUsecase(authRepo, passSvc, jwtAdmin, "admin")
	patientAuthUC := _authUsecases.NewAuthUsecase(authRepo, passSvc, jwtPatient, "patient")

	// Auth endpoints
	authGroup := v1.Group("/auth")
	_authControllers.NewAuthController(authGroup.Group("/admin"), adminAuthUC)
	_authControllers.NewAuthController(authGroup.Group("/patient"), patientAuthUC)

	// Role-based JWT middleware
	adminAuth := utils.NewJWTMiddleware(jwtAdmin, "admin")
	// patientAuth := utils.NewJWTMiddleware(jwtPatient, "patient")

	// ------------------ ADMIN MODULES ------------------
	adminGroup := v1.Group("/admins", adminAuth)

	adminRepo := _adminRepositories.NewAdminRepository(s.Db)
	adminUsecase := _adminUsecases.NewAdminUsecase(adminRepo, passSvc)
	_adminControllers.NewAdminController(adminGroup, adminUsecase)

	diseaseRepo := _adminRepositories.NewDiseaseGetRepository(s.Db)
	diseaseUC := _adminUsecases.NewDiseaseUsecase(diseaseRepo)
	_adminControllers.NewDiseaseController(adminGroup.Group("/diseases"), diseaseUC)

	tx := _adminRepositories.NewTransactionGorm(s.Db)
	patientUC := _adminUsecases.NewPatientUsecase(tx, passSvc)
	patientGetRepo := _adminRepositories.NewPatientGetRepository(s.Db)
	patientGetUC := _adminUsecases.NewPatientGetUsecase(patientGetRepo)
	_adminControllers.NewPatientController(adminGroup.Group("/patients"), patientUC, patientGetUC, patientUC)

	// ------------------ PATIENT MODULES ------------------

	// ------------------ 404 HANDLER ------------------
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
