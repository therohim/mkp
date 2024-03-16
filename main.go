package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"test-mkp/config"
	"test-mkp/exception"
	"test-mkp/middleware"
	"test-mkp/utils"

	user_ctrl "test-mkp/src/user/controller"
	user_repo "test-mkp/src/user/repository"
	user_svc "test-mkp/src/user/service"

	ticket_ctrl "test-mkp/src/ticket/controller"
	ticket_repo "test-mkp/src/ticket/repository"
	ticket_svc "test-mkp/src/ticket/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	configuration := config.New()
	jwtConfig := config.NewJwtConfig(configuration)
	postgreDb := config.NewPostgreSqlxDatabase(configuration)
	authMiddleware := middleware.NewAuthMiddleware(postgreDb.DB, jwtConfig)

	// Initialize pkg
	utils.NewLogger()

	// setup repository
	userRepo := user_repo.NewUserRepository(postgreDb.DB)
	ticketRepo := ticket_repo.NewTicketRepository(postgreDb.DB)

	// setup service
	userService := user_svc.NewUserService(userRepo, jwtConfig)
	ticketService := ticket_svc.NewTicketService(ticketRepo)

	// setup controller
	userCtrl := user_ctrl.NewUserController(userService)
	ticketCtrl := ticket_ctrl.NewTicketController(ticketService, *authMiddleware)

	// Setup Fiber
	app := fiber.New(config.NewFiberConfig())
	app.Use(recover.New())
	// app.Use(cors.New(config.NewCorsConfig()))
	app.Use(config.NewLogger(config.LoggerConfig{
		Logger: utils.Logger,
	}))

	// Setup Routing
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(map[string]interface{}{
			"code":   200,
			"status": "success",
			"IP":     ctx.IP(),
			"IPs":    ctx.IPs(),
		})
	})

	//Register Rest API Route
	userCtrl.Route(app)
	ticketCtrl.Route(app)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start App
	go func() {
		if err := app.Listen(":3000"); err != nil {
			exception.PanicIfNeeded(err)
		}
	}()

	// graceful shutdown
	<-stop
	log.Println("Stopping server...")
	postgreDb.DB.Close()
}
