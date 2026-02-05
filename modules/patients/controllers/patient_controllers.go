package controllers

import (
	"strconv"
	"strings"

	"github.com/chothanin01/PhoSS-Care-server/modules/patients/entities"
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
}

func (c *PatientController) CreatePatient(ctx *fiber.Ctx) error {
	req := new(entities.PatientFullCreateReq)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := c.PatientUsecase.CreateFull(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"message":    "User, Patient, and Relatives created successfully",
		"data":       res,
		"statusCode": fiber.StatusOK,
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

