package controllers

import (
	"main/app/services"
	
	"github.com/gofiber/fiber/v2"
)


// @Router /v1/s3/upload/{email} [post]
func UploadHandler(c *fiber.Ctx) error {
	email := c.Params("email")

	err := services.UploadToS3(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
	})
}

// @Router /v1/s3/download/{email} [get]
func DownloadHandler(c *fiber.Ctx) error {
	email := c.Params("email")

	zipFilePath, err := services.DownloadToZip(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).SendFile(zipFilePath)
}