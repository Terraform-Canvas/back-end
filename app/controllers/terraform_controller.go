package controllers

import (
	"path/filepath"

	"main/app/services"

	"github.com/gofiber/fiber/v2"
)

// Create a user env tf
// @Router /v1/terraform/merge/{email} [post]
func MergeEnvTf(c *fiber.Ctx) error {
	var data []map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	email := c.Params("email")
	userFolderPath := filepath.Join("usertf", email)

	err := services.InitializeFolder(userFolderPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	err = services.MergeEnvTf(userFolderPath, data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	err = services.CreateTfvars(userFolderPath, data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "merge user env tf success",
	})
}

// Apply user env's tf
// @Router /v1/terraform/apply/{email} [post]
func ApplyEnvTf(c *fiber.Ctx) error {
	email := c.Params("email")
	result, err := services.ApplyTerraform(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   result,
	})
}
