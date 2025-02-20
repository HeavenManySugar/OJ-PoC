package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/HeavenManySugar/OJ-PoC/docs"
	"github.com/HeavenManySugar/OJ-PoC/handlers"
)

// New create an instance of Book app routes
func New() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "UTC",
	}))
	
	app.Get("/swagger/*", swagger.HandlerDefault) 
	api := app.Group("api")


	api.Get("/books", handlers.GetAllBooks)
	api.Get("/books/:id", handlers.GetBookByID)
	api.Post("/books", handlers.RegisterBook)
	api.Delete("/books/:id", handlers.DeleteBook)
	api.Post("/gitea", handlers.PostGiteaHook)
	api.Post("/sandbox", handlers.PostSandboxCmd)
	api.Get("/scores", handlers.GetScores)
	api.Get("/score", handlers.GetScoreByRepo)

	return app
}