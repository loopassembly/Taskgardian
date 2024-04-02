// auth.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"GDSC/controllers"
	"GDSC/middleware"
)
// all auth routes including oauth
func SetupAuthRoutes(router fiber.Router) {
	router.Post("/register", controllers.AdminSignIn)
	router.Post("/login", controllers.SignInUser)
	router.Get("/logout", middleware.DeserializeUser, controllers.LogoutUser)
	router.Get("/verifyemail/:verificationCode", controllers.VerifyEmail)
	router.Post("/forgotpassword", controllers.ForgotPassword)
	router.Patch("/resetpassword/:resetToken",controllers.ResetPassword)
	
	router.Get("/getinfo/:id",controllers.GetUserTasks)
	router.Get("/task",controllers.CreateTask)
	// router.Get("/sessions/oauth/google", controllers.GoogleOAuth)
	// router.Get("/sessions/oauth/github", controllers.GitHubOAuth)
}