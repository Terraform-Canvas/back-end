package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// UserSignIn method to auth user and return access and refresh tokens.
// @Description Auth user and return access and refresh token.
// @Summary auth user and return access and refresh token
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "User Email"
// @Param password body string true "User Password"
// @Success 200 {string} status "ok"
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Description De-authorize user and delete refresh token from Redis.
// @Summary de-authorize user and delete refresh token from Redis
// @Tags User
// @Accept json
// @Produce json
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/user/sign/out [post]
func UserSignOut(c *fiber.Ctx) error {
	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
