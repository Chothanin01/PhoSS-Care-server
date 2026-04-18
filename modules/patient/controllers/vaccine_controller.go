package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type VaccineQueryController struct {
	usecase entities.VaccineQueryUsecase
}

func NewVaccineQueryController(r fiber.Router, uc entities.VaccineQueryUsecase) {
	controller := &VaccineQueryController{usecase: uc}
	r.Get("/", controller.GetVaccines)
	r.Get("/:id", controller.GetVaccineDetail)
}

func (c *VaccineQueryController) GetVaccines(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "unauthorized"})
	}
	
	patientID, err := uuid.Parse(claims["role_id"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "invalid token"})
	}

	page := ctx.QueryInt("page", 1)         
	filter := ctx.Query("filter", "all")    

	response, err := c.usecase.GetVaccineList(patientID, page, filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

func (c *VaccineQueryController) GetVaccineDetail(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "unauthorized"})
	}
	
	patientID, err := uuid.Parse(claims["role_id"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "invalid token"})
	}

	vaccineID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid vaccine ID format"})
	}

	detail, err := c.usecase.GetVaccineDetail(patientID, vaccineID)
	if err != nil {
		
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "vaccine not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    detail,
	})
}