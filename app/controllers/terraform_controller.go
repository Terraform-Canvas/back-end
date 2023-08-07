package controllers

import (
	"path/filepath"
	"time"

	"main/pkg/utils"

	"main/app/models"
	"main/app/services"

	"github.com/gofiber/fiber/v2"
)

// Create a user env tf
// @Router /v1/terraform/usertf [post]
func MergeEnvTf(c *fiber.Ctx) error {
	now := time.Now().Unix()
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	expires := claims.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	resources := []models.Resource{}
	if err := c.BodyParser(&resources); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	email, err := utils.GetEmailFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	userFolderPath := filepath.Join("usertf", email)

	err = services.InitializeFolder(userFolderPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	err = services.MergeEnvTf(userFolderPath, resources)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	err = services.CreateTfvars(userFolderPath, resources)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	err = services.ApplyTerraform(userFolderPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "Successful creation of the user's tf files environment",
	})
}

// @Router /v1/terraform/destroy [post]
func DestroyEnv(c *fiber.Ctx) error {
	now := time.Now().Unix()
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	expires := claims.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	email, err := utils.GetEmailFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	userFolderPath := filepath.Join("usertf", email)

	err = services.DestroyTerraform(userFolderPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "successful env destruction",
	})
}
