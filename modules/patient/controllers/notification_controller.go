package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	
	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
)

type NotificationController struct {
	queryUsecase   entities.NotificationQueryUsecase
	commandUsecase entities.NotificationCommandUsecase
}

func NewNotificationController(r fiber.Router, qUc entities.NotificationQueryUsecase, cUc entities.NotificationCommandUsecase) {
	controller := &NotificationController{queryUsecase: qUc, commandUsecase: cUc}

	r.Get("/", controller.GetAllNotifications)           
	r.Get("/status", controller.CheckUnread)      
	r.Patch("/:id/read", controller.MarkNotificationRead)
}

func getPatientID(ctx *fiber.Ctx) (uuid.UUID, error) {
	claims := ctx.Locals("user").(jwt.MapClaims)
	return uuid.Parse(claims["role_id"].(string))
}

func (c *NotificationController) GetAllNotifications(ctx *fiber.Ctx) error {
	patientID, _ := getPatientID(ctx)
	
	notis, err := c.queryUsecase.GetPatientNotifications(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return ctx.JSON(fiber.Map{"success": true, "data": notis})
}

func (c *NotificationController) CheckUnread(ctx *fiber.Ctx) error {
	patientID, _ := getPatientID(ctx)

	hasUnread, err := c.queryUsecase.CheckHasUnread(patientID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return ctx.JSON(fiber.Map{"success": true, "has_unread": hasUnread})
}

func (c *NotificationController) MarkNotificationRead(ctx *fiber.Ctx) error {
	patientID, _ := getPatientID(ctx)

	notiID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid notification ID"})
	}

	err = c.commandUsecase.MarkAsRead(notiID, patientID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "notification not found"})
	}

	return ctx.JSON(fiber.Map{"success": true, "message": "marked as read"})
}