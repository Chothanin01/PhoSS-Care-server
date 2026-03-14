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
	PatientUpdateUsecase entities.PatientUpdateUsecase
}

func NewPatientController(r fiber.Router, createUC entities.PatientUsecase, getUC entities.PatientGetUsecase, updateUC entities.PatientUpdateUsecase) {
	controller := &PatientController{
		PatientUsecase:    createUC,
		PatientGetUsecase: getUC,
		PatientUpdateUsecase: updateUC,
	}	

	r.Post("/", controller.CreatePatient)
	r.Get("/", controller.GetPatients)
	r.Get("/:id", controller.GetPatientByID)
	r.Get("/:id/info", controller.GetPatientBasicInfo)
	r.Get("/:id/appointments", controller.GetPatientAppointmentsInfo)
	r.Get("/:id/diseases" , controller.GetPatientDiseases)
	r.Get("/:id/vaccines", controller.GetPatientVaccines)
	r.Get("/:id/:disease_id", controller.GetPatientDiseaseInfo)
	r.Patch("/:id", controller.UpdatePatient)

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

	res, err := c.PatientUsecase.CreateFull(req, &creatorID)
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

func (c *PatientController) GetPatientAppointmentsInfo(ctx *fiber.Ctx) error {
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

	res, err := c.PatientGetUsecase.GetPatientAppointmentsInfoByID(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *PatientController) GetPatientVaccines(ctx *fiber.Ctx) error {
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

	res, err := c.PatientGetUsecase.GetPatientVaccinesByID(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (c *PatientController) GetPatientBasicInfo(ctx *fiber.Ctx) error {

	idParam := ctx.Params("id")

	patientID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid patient id",
		})
	}

	data, err := c.PatientGetUsecase.GetPatientBasicInfoByID(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (c *PatientController) GetPatientDiseases(ctx *fiber.Ctx) error {

	patientID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid patient id",
		})
	}

	dtype := ctx.Query("type", "active")

	diseases, err := c.PatientGetUsecase.GetPatientDiseases(patientID, dtype)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Get patient diseases successfully.",
		"data": diseases,
	})
}

func (c *PatientController) UpdatePatient(ctx *fiber.Ctx) error {
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

	var req entities.PatientUpdateReq

	claims := ctx.Locals("user").(jwt.MapClaims)
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	adminID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id"})
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	res, err := c.PatientUpdateUsecase.UpdatePatientInfo(patientID, &req, adminID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Patient updated successfully",
		"data":    res,
	})
}

