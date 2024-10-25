package routes

import (
	"go-rest-app/pkg/controllers"

	"github.com/gofiber/fiber/v2"
)

func CardRoutes(app *fiber.App) {
	app.Get("/cards", controllers.GetCards)
	app.Get("/cards/:id", controllers.GetCard)
	app.Post("/card", controllers.CreateCard)
	app.Put("/cards/:id", controllers.UpdateCard)
	app.Delete("/cards/:id", controllers.DeleteCard)
}
