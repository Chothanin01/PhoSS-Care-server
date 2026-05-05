package controllers

import (

	"github.com/chothanin01/PhoSS-Care-server/modules/patient/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PatientController struct {
	GetUsecase entities.GetPatientUsecase
}

func NewPatientController(r fiber.Router, uc entities.GetPatientUsecase) {
	controller := &PatientController{GetUsecase: uc}
	r.Get("/basicinfo", controller.GetPatientBasicInfo)
	r.Get("/fullinfo", controller.GetPatientFullInfo)
	r.Get("/diseases", controller.GetPatientDiseases)
}

func (c *PatientController) GetPatientBasicInfo(ctx *fiber.Ctx) error {
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

	data, err := c.GetUsecase.GetPatientBasicInfo(patientID)
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

func (c *PatientController) GetPatientFullInfo(ctx *fiber.Ctx) error {
    claims, ok := ctx.Locals("user").(jwt.MapClaims)
    if !ok {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "message": "unauthorized access",
        })
    }

    userIDStr, ok := claims["user_id"].(string)
    if !ok {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "message": "invalid token claims: user_id missing",
        })
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "invalid user format",
        })
    }

    data, err := c.GetUsecase.GetPatientFullInfo(userID)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    // 4. Response
    return ctx.JSON(fiber.Map{
        "success": true,
        "message": "Get patient full info successfully.",
        "data":    data,
    })
}


func (c *PatientController) GetPatientDiseases(ctx *fiber.Ctx) error {
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, 
			"message": "unauthorized",
		})
	}
	
	patientIDStr, ok := claims["role_id"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, 
			"message": "invalid token payload",
		})
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, 
			"message": "invalid token payload",
		})
	}

	filter := entities.DiseaseFilter{
		Type: ctx.Query("type"),
	}

	diseases, err := c.GetUsecase.GetPatientDiseases(patientID, filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    diseases,
	})
}