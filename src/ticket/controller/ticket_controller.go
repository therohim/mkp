package controller

import (
	"net/http"
	"test-mkp/exception"
	"test-mkp/middleware"
	"test-mkp/src/ticket/model"
	"test-mkp/src/ticket/service"
	"test-mkp/utils"

	"github.com/gofiber/fiber/v2"
)

type TicketController struct {
	ticketService service.TicketService
	middleware    middleware.AuthMiddleware
}

func NewTicketController(ticketService service.TicketService, middleware middleware.AuthMiddleware) *TicketController {
	return &TicketController{ticketService, middleware}
}

func (controller *TicketController) Route(app *fiber.App) {
	v1 := app.Group("/api/v1")
	ticket := v1.Group("/ticket")
	ticket.Get("/", controller.middleware.AuthMiddleware, controller.getTicket)
	ticket.Post("/", controller.middleware.AuthMiddleware, controller.add)
	ticket.Put("/:ticketId", controller.middleware.AuthMiddleware, controller.editTicket)
	ticket.Delete("/:ticketId", controller.middleware.AuthMiddleware, controller.deleteTicket)
}

func (controller *TicketController) getTicket(ctx *fiber.Ctx) error {
	req := model.ListRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		exception.PanicAndLog(exception.ClientError{
			Message: err.Error(),
		})
	}

	results, pagination, err := controller.ticketService.GetTicket(req)
	if err != nil {
		exception.PanicAndLog(exception.ClientError{
			Message: err.Error(),
		})
	}

	return ctx.JSON(utils.WebArrayResponse{
		Code:     http.StatusOK,
		Status:   "success",
		Data:     results,
		Paginate: pagination,
	})
}

func (controller *TicketController) add(c *fiber.Ctx) error {
	var request model.AddTicketRequest
	if err := c.BodyParser(&request); err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	err := controller.ticketService.AddTicket(request)
	if err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	return c.JSON(utils.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Ticket successfull addedd",
	})
}

func (controller *TicketController) deleteTicket(c *fiber.Ctx) error {
	ticketID := c.Params("ticketId")

	err := controller.ticketService.DeleteTicket(ticketID)
	if err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	return c.JSON(utils.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Ticket successfull deleted",
	})
}

func (controller *TicketController) editTicket(c *fiber.Ctx) error {
	ticketID := c.Params("ticketId")

	var request model.EditTicketRequest
	if err := c.BodyParser(&request); err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}
	err := controller.ticketService.EditTicket(ticketID, request)
	if err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	return c.JSON(utils.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Ticket successfull updated",
	})
}
