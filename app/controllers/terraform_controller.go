package controllers

import (
	"main/app/services"
	"github.com/gofiber/fiber/v2"
)

//accept랑 produce는 (swagger에서의) 애매해서 생략(s3에서도 생략되어 있음)
//Description Merge module tf file and Add in user folder
//Summary merge module tf file and add in user folder
//Tags Env
//Param email path string true "user email"
//Success 200 {string} status "ok"

// @Router /v1/terraform/merge/{email} [post]
func MergeEnvTf(c *fiber.Ctx) error {
	var data []map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	email := c.Params("email")

	err := services.MergeEnvTf(email, data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
	})
}

// @Router /v1/terraform/tfvars/{email} [post]
func CreateTfvars(c *fiber.Ctx) error {
	var data []map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	email := c.Params("email")

	err := services.CreateTfvars(email, data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
	})
}

// @Router /v1/terraform/apply/{email} [post]
func ApplyEnvTf(c *fiber.Ctx) error {
	email := c.Params("email")
	result, err := services.ApplyTerraform(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg": result,
	})
}