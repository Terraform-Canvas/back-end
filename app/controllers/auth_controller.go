package controllers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"main/app/models"
)

// UserSignIn method to auth user and return access and refresh tokens.
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	// Return status 200 OK.
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println(user)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.JSON(fiber.Map{
		"authInfo": map[string]string{
			"accessToken":  "testAccess",
			"refreshToken": "testRefresh",
		},
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Router /v1/user/sign/out [post]
func UserSignOut(c *fiber.Ctx) error {
	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
