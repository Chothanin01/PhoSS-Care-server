package controllers

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type RequestController struct {
	usecase entities.RequestGetUsecase
}

func NewRequestController(r fiber.Router, uc entities.RequestGetUsecase) {
	controller := &RequestController{usecase: uc}
	r.Get("/", controller.GetRequests)
	r.Get("/:id", controller.GetRequestDetailByID)
}

func (c *RequestController) GetRequests(ctx *fiber.Ctx) error {
	reqType := ctx.Query("req_type")
	status := ctx.Query("status")

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	params := entities.RequestQueryParams{
		ReqType: reqType,
		Status:  status,
		Page:    page,
		Limit:   limit,
	}

	res, err := c.usecase.GetRequestsWithFilter(params)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Requests fetched successfully",
		"data":    res,
	})
}

func (c *RequestController) GetRequestDetailByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request id",
		})
	}

	res, err := c.usecase.GetRequestInfoByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Request detail fetched successfully",
		"data":    res,
	})
}