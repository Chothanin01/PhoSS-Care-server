package controllers

import (
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
	"github.com/gofiber/fiber/v2"
)

type PatientController struct {
	PatientUsecase    entities.PatientUsecase
	PatientGetUsecase entities.PatientGetUsecase
}

func NewPatientController(r fiber.Router, createUC entities.PatientUsecase, getUC entities.PatientGetUsecase) {
	controller := &PatientController{
		PatientUsecase:    createUC,
		PatientGetUsecase: getUC,
	}	

	r.Post("/", controller.CreatePatient)
	r.Get("/", controller.GetPatients)
	r.Get("/:id", controller.GetPatientByID)
	r.Get("/:id/appointments", controller.GetPatientAppointments)
	r.Get("/:id/:disease_id", controller.GetPatientDiseaseInfo)

}

func (c *PatientController) CreatePatient(ctx *fiber.Ctx) error {
	req := new(entities.PatientFullCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	claims := ctx.Locals("user").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	creatorID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id"})
	}

	res, err := c.PatientUsecase.CreateFull(req, creatorID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Patient created successfully",
		"data":    res,
	})
}

func (c *PatientController) GetPatients(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	search := ctx.Query("search", "")

	queryMap := ctx.Queries()
	diseases := []string{}
	for key, val := range queryMap {
		if key == "disease" {
			parts := strings.Split(val, ",")
			for _, p := range parts {
				if p != "" {
					diseases = append(diseases, strings.TrimSpace(p))
				}
			}
		}
	}

	appointStr := ctx.Query("appoint", "")
	var appoint *bool
	if appointStr != "" {
		val := appointStr == "true"
		appoint = &val
	}

	req := entities.PatientQueryParams{
		Search:   search,
		Diseases: diseases,
		Appoint:  appoint,
		Page:     page,
		Limit:    limit,
	}

	res, err := c.PatientGetUsecase.GetPatientListWithFilter(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *PatientController) GetPatientByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Missing patient ID",
		})
	}

	patientID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid patient ID format",
		})
	}

	res, err := c.PatientGetUsecase.GetPatientInfoByID(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *PatientController) GetPatientDiseaseInfo(ctx *fiber.Ctx) error {
	patientIDParam := ctx.Params("id")
	diseaseIDParam := ctx.Params("disease_id")

	patientID, err := uuid.Parse(patientIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid patient ID",
		})
	}

	diseaseID, err := uuid.Parse(diseaseIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid disease ID",
		})
	}

	data, err := c.PatientGetUsecase.GetPatientDiseasesInfo(patientID, diseaseID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Get patient disease info successfully.",
		"data":    data,
	})
}

func (c *PatientController) GetPatientAppointments(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Missing patient ID",
		})
	}

	patientID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid patient ID format",
		})
	}

	res, err := c.PatientGetUsecase.GetPatientAppointmentsByID(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
