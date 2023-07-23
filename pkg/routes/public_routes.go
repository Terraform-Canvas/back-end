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
	route.Post("/terraform/usertf", controllers.MergeEnvTf)
	route.Post("/terraform/apply", controllers.ApplyEnvTf)
	route.Post("/s3/upload", controllers.UploadHandler)
	route.Get("/s3/download", controllers.DownloadHandler)
	route.Post("/user/key", controllers.UserKeySave)
	route.Get("/ec2/instanceTypes", controllers.EC2InstanceTypes)
	route.Get("/ec2/ami", controllers.EC2Images)
}
