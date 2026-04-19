package controllers

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AppointmentController struct {
	QueryUC entities.AppointmentQueryUsecase
	CommandUC entities.AppointmentCommandUsecase
}

func NewAppointmentController(r fiber.Router, Queryuc entities.AppointmentQueryUsecase, Commanduc entities.AppointmentCommandUsecase) {
	controller := &AppointmentController{
		QueryUC: Queryuc,
		CommandUC: Commanduc,
	}
	
	r.Get("/", controller.ListPatientAppointments)
	r.Get("/history/:disease_id", controller.GetDiseaseHistory)
	r.Get("/history/detail/:appoint_id", controller.GetHistoryDetail)
	r.Get("/schedule/:disease_id", controller.GetDiseaseSchedule)
	r.Get("/:disease_id", controller.GetAppointmentDetail)
	r.Post("/delay", controller.CreateDelayRequest)
	r.Patch("/:appoint_id/canceldelay", controller.CancelDelayRequest)
	}

func (c *AppointmentController) GetAppointmentDetail(ctx *fiber.Ctx) error {

	claims := ctx.Locals("user").(jwt.MapClaims)
	patientIDStr := claims["role_id"].(string)
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}

	diseaseIDStr := ctx.Params("disease_id")
	diseaseID, err := uuid.Parse(diseaseIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid disease id"})
	}

	data, err := c.QueryUC.GetAppointmentDetail(patientID, diseaseID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "appointment not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (c *AppointmentController) ListPatientAppointments(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized access",
		})
	}

	patientIDStr, ok := claims["role_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid token claims",
		})
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	data, err := c.QueryUC.ListPatientAppointments(patientID)
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

func (c *AppointmentController) CreateDelayRequest(ctx *fiber.Ctx) error {

	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	userID, err1 := uuid.Parse(claims["user_id"].(string))
	patientID, err2 := uuid.Parse(claims["role_id"].(string))

	if err1 != nil || err2 != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token data"})
	}

	payload := new(entities.AppointmentDelayReq)
	if err := ctx.BodyParser(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	err := c.CommandUC.SubmitDelayRequest(userID, patientID, payload)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "delay request submitted successfully",
	})
}

func (c *AppointmentController) GetDiseaseSchedule(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized access"})
	}

	patientIDStr, ok := claims["role_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token claims"})
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid patient ID in token"})
	}

	diseaseIDStr := ctx.Params("disease_id")
	diseaseID, err := uuid.Parse(diseaseIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid disease ID format",
		})
	}

	data, err := c.QueryUC.GetScheduleByDisease(patientID, diseaseID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (c *AppointmentController) GetDiseaseHistory(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}
	patientID, _ := uuid.Parse(claims["role_id"].(string))

	diseaseID, err := uuid.Parse(ctx.Params("disease_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid disease ID"})
	}

	page := ctx.QueryInt("page", 1)

	data, err := c.QueryUC.GetDiseaseHistory(patientID, diseaseID, page)
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

func (c *AppointmentController) GetHistoryDetail(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}
	patientID, _ := uuid.Parse(claims["role_id"].(string))

	appointIDStr := ctx.Params("appoint_id")
	appointID, err := uuid.Parse(appointIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid appointment ID format",
		})
	}

	data, err := c.QueryUC.GetHistoryDetail(appointID, patientID)
	if err != nil {
		if err.Error() == "failed to fetch appointment details: appointment not found or access denied" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "appointment not found",
			})
		}

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

func (c *AppointmentController) CancelDelayRequest(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "unauthorized"})
	}
	
	patientID, err := uuid.Parse(claims["role_id"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "invalid token"})
	}

	appointID, err := uuid.Parse(ctx.Params("appoint_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid appointment ID format"})
	}

	err = c.CommandUC.CancelDelayRequest(appointID, patientID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// 4. Return Success
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Delay request has been canceled successfully",
	})
}