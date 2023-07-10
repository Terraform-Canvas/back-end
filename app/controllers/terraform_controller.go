package controllers

import (
	"main/app/services"
	"github.com/gofiber/fiber/v2"
)

//accept랑 produce는 (swagger에서의) 애매해서 생략(s3에서도 생략되어 있음)
//Description Merge module tf file and Add in user folder
//Summary merge module tf file and add in user folder
//Tags Env
//Param usermail path string true "user mail"
//Success 200 {string} status "ok"

// @Router /v1/terraform/merge/{usermail} [post]
func MergeEnvTf(c *fiber.Ctx) error {
	var data []map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	usermail := c.Params("usermail")

	err := services.MergeEnvTf(usermail, data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	return c.Status(fiber.StatusOK).SendString(usermail)
}

// @Router /v1/terraform/tfvars/{usermail} [post]
func CreateTfvars(c *fiber.Ctx) error {
	var data []map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	usermail := c.Params("usermail")

	err := services.CreateTfvars(usermail, data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	return c.Status(fiber.StatusOK).SendString(usermail)
}

// @Router /v1/terraform/apply/{usermail} [post]
func ApplyEnvTf(c *fiber.Ctx) error {
	usermail := c.Params("usermail")
	result, err := services.ApplyTerraform(usermail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": result,
	})
}