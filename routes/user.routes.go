// user.go
package routes

import (
	"GDSC/controllers"
	"GDSC/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	router.Get("/me", middleware.DeserializeUser, controllers.GetMe)
}
