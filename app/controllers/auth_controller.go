package controllers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"main/app/models"
	"main/platform/database"
)

// UserSignIn method to auth user and return access and refresh tokens.
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	// Return status 200 OK.
	signIn := &models.SignIn{}

	if err := c.BodyParser(signIn); err != nil {
		log.Println(signIn)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	db, err := database.OCINoSQLConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	user, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given email is not found",
		})
	}
	if user.Password != signIn.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given password is not correct",
		})
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
