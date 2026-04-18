package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type requestController struct {
	CommandUC entities.RequestCommandUsecase
	QueryUC   entities.RequestQueryUsecase
}

func NewRequestCommandController(r fiber.Router, CommandUC entities.RequestCommandUsecase, QueryUC entities.RequestQueryUsecase) {
	controller := &requestController{
		CommandUC: CommandUC,
		QueryUC:   QueryUC,
	}

	r.Get("/", controller.GetOptions)
	r.Post("/", controller.CreateDocumentRequests)
}

func (c *requestController) CreateDocumentRequests(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	userID, _ := uuid.Parse(claims["user_id"].(string))
	patientID, _ := uuid.Parse(claims["role_id"].(string))

	payload := new(entities.CreateDocumentReq)
	if err := ctx.BodyParser(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	err := c.CommandUC.SubmitDocumentRequest(userID, patientID, payload)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "requests submitted successfully",
	})
}

func (c *requestController) GetOptions(ctx *fiber.Ctx) error {

	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}
	
	patientID, err := uuid.Parse(claims["role_id"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
	}

	options, err := c.QueryUC.GetAvailableDocumentOptions(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to fetch options",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    options,
	})
}