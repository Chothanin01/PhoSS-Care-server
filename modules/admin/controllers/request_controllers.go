package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/chothanin01/PhoSS-Care-server/modules/admin/entities"
)

type RequestController struct {
	usecase entities.RequestGetUsecase
}

func NewRequestController(r fiber.Router, uc entities.RequestGetUsecase) {
	controller := &RequestController{usecase: uc}
	r.Get("/", controller.GetRequests)
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