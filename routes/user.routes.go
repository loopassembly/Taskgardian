// user.go
package routes

import (
	"GDSC/controllers"
	"GDSC/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	router.Post("/updatetask/:taskid",middleware.DeserializeUser,controllers.UpdateTask)
	router.Post("/task",middleware.DeserializeUser,controllers.CreateTask)
	router.Get("/me", middleware.DeserializeUser, controllers.GetMe)
	router.Delete("/DeleteTask/:id", middleware.DeserializeUser, controllers.DeleteTask)
	router.Get("/updateUserRole", middleware.DeserializeUser, controllers.UpdateUserRole)
}
