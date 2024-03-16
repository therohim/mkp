package controller

import (
	"fmt"
	"net/http"
	"test-mkp/exception"
	"test-mkp/src/user/model"
	"test-mkp/src/user/service"
	"test-mkp/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (controller *UserController) Route(app *fiber.App) {
	v1 := app.Group("/api/v1")
	auth := v1.Group("/auth")
	auth.Post("/register", controller.register)
	auth.Post("/login", controller.login)
}

func (controller *UserController) register(c *fiber.Ctx) error {
	fmt.Println(c.Locals("userId"))
	var request model.RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	err := controller.userService.Register(request)
	if err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	return c.JSON(utils.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Register successfull",
	})
}

func (controller *UserController) login(c *fiber.Ctx) error {
	var request model.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	response, err := controller.userService.Login(request)
	if err != nil {
		exception.PanicIfNeeded(exception.ClientError{
			Message: err.Error(),
		})
	}

	return c.JSON(utils.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Login successfull",
		Data:    response,
	})
}
