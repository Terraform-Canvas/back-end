package routes

import (
	"github.com/gofiber/fiber/v2"
	"main/app/controllers"
	"main/pkg/middleware"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/login/refresh", middleware.JWTProtected(), controllers.UserRefresh)
	route.Post("/logout", middleware.JWTProtected(), controllers.UserSignOut)
	route.Post("/terraform/usertf", middleware.JWTProtected(), controllers.MergeEnvTf)
	route.Post("/terraform/apply", middleware.JWTProtected(), controllers.ApplyEnvTf)
	route.Post("/s3/upload", middleware.JWTProtected(), controllers.UploadHandler)
	route.Get("/s3/download", middleware.JWTProtected(), controllers.DownloadHandler)
	route.Post("/user/key", middleware.JWTProtected(), controllers.UserKeySave)
}
