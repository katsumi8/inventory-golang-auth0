package order

import "github.com/gofiber/fiber/v2"

func AddOrderRoutes(app *fiber.App, controller *OrderController) {
	todos := app.Group("/orders")

	// add middlewares here

	// add routes here
	todos.Post("/", controller.create)
	todos.Get("/", controller.getAll)
}
