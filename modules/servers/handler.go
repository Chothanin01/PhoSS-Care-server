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

	_patientControllers "github.com/chothanin01/PhoSS-Care-server/modules/patient/controllers"
	_patientRepositories "github.com/chothanin01/PhoSS-Care-server/modules/patient/repositories"
	_patientUsecases "github.com/chothanin01/PhoSS-Care-server/modules/patient/usecases"

	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	// ------------------ AUTH SETUP ------------------
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
	patientAuth := utils.NewJWTMiddleware(jwtPatient, "patient")

	// ------------------ ADMIN MODULES ------------------
	adminGroup := v1.Group("/admins", adminAuth)

	adminRepo := _adminRepositories.NewAdminRepository(s.Db)
	adminUsecase := _adminUsecases.NewAdminUsecase(adminRepo, passSvc)
	_adminControllers.NewAdminController(adminGroup, adminUsecase)

	adminDiseaseRepo := _adminRepositories.NewDiseaseGetRepository(s.Db)
	adminDiseaseUC := _adminUsecases.NewDiseaseUsecase(adminDiseaseRepo)
	_adminControllers.NewDiseaseController(adminGroup.Group("/diseases"), adminDiseaseUC)

	adminRelativeRepo := _adminRepositories.NewRelativeRepository(s.Db)
	adminRelativeUC := _adminUsecases.NewRelativeUsecase(adminRelativeRepo)
	_adminControllers.NewRelativeController(adminGroup.Group("/patients"), adminRelativeUC)

	tx := _adminRepositories.NewTransactionGorm(s.Db)
	adminPatientUC := _adminUsecases.NewPatientUsecase(tx, passSvc)
	adminPatientGetRepo := _adminRepositories.NewPatientGetRepository(s.Db)
	adminPatientGetUC := _adminUsecases.NewPatientGetUsecase(adminPatientGetRepo)
	_adminControllers.NewPatientController(adminGroup.Group("/patients"), adminPatientUC, adminPatientGetUC, adminPatientUC)

	adminAppointTx := _adminRepositories.NewTransactionGorm(s.Db)
	adminAppointRepo := _adminRepositories.NewAppointmentRepository(s.Db)
	adminAppointUC := _adminUsecases.NewAppointmentUsecase(adminAppointTx, adminAppointRepo)
	_adminControllers.NewAppointmentController(adminGroup.Group("/appointments"), adminAppointUC)

	adminRequestRepo := _adminRepositories.NewRequestGetRepository(s.Db)
	adminRequestUC := _adminUsecases.NewRequestGetUsecase(adminRequestRepo)
	adminRequestUpdateRepo := _adminRepositories.NewRequestUpdateRepository(s.Db)
	adminRequestUpdateUC := _adminUsecases.NewRequestUpdateUsecase(adminRequestUpdateRepo)
	_adminControllers.NewRequestController(adminGroup.Group("/requests"), adminRequestUC, adminRequestUpdateUC)

	// ------------------ PATIENT MODULES ------------------
	patientGroup := v1.Group("patient", patientAuth)

	patientGetRepo := _patientRepositories.NewGetPatientRepository(s.Db)
	patientGetUC := _patientUsecases.NewGetPatientUsecase(patientGetRepo)
	_patientControllers.NewPatientController(patientGroup, patientGetUC)

	patientQueryAppointRepo := _patientRepositories.NewappointmentQueryRepo(s.Db)
	patientQueryAppointUC := _patientUsecases.NewAppointmentQueryUsecase(patientQueryAppointRepo)

	patientAppointCommandRepo := _patientRepositories.NewAppointmentCommandRepo(s.Db)
	patientAppointCommandUC := _patientUsecases.NewAppointmentCommandUsecase(patientAppointCommandRepo)
	
	_patientControllers.NewAppointmentController(patientGroup.Group("/appointments"), patientQueryAppointUC, patientAppointCommandUC)

	patientRequestCommandRepo := _patientRepositories.NewRequestCommandRepo(s.Db)
	patientRequestCommandUC := _patientUsecases.NewRequestCommandUsecase(patientRequestCommandRepo)

	patientRequestQueryRepo := _patientRepositories.NewRequestQueryRepo(s.Db)
	patientRequestQueryUC := _patientUsecases.NewRequestQueryUsecase(patientRequestQueryRepo)

	_patientControllers.NewRequestCommandController(patientGroup.Group("/requests"), patientRequestCommandUC, patientRequestQueryUC)

	patientNotiCommandRepo := _patientRepositories.NewNotificationCommandRepo(s.Db)
	patientNotiCommandUC := _patientUsecases.NewNotificationCommandUsecase(patientNotiCommandRepo)
	
	patientNotiQueryRepo := _patientRepositories.NewNotificationQueryRepo(s.Db)
	patientNotiQueryUC := _patientUsecases.NewNotificationQueryUsecase(patientNotiQueryRepo)

	_patientControllers.NewNotificationController(patientGroup.Group("/noti"), patientNotiQueryUC, patientNotiCommandUC)

	patientVaccineQueryRepo := _patientRepositories.NewVaccineQueryRepo(s.Db)
	patientVaccineQueryUC := _patientUsecases.NewVaccineQueryUsecase(patientVaccineQueryRepo)

	_patientControllers.NewVaccineQueryController(patientGroup.Group("/vaccine"), patientVaccineQueryUC)
	
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
