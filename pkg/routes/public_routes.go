package routes

import (
	"github.com/gofiber/fiber/v2"
	"main/app/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")
	// Routes for POST method:
	route.Post("/login/new", controllers.UserSignIn)
	route.Post("/terraform/usertf/:email", controllers.MergeEnvTf)
	route.Post("/terraform/apply/:email", controllers.ApplyEnvTf)
	route.Post("/s3/upload/:email", controllers.UploadHandler)
	route.Get("/s3/download/:email", controllers.DownloadHandler)
	route.Post("/user/key", controllers.UserKeySave)
}
