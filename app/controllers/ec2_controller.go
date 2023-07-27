package controllers

import (
	"github.com/gofiber/fiber/v2"
	"main/pkg/utils"
	"main/platform/amazon"
	"time"
)

// @Router /v1/ec2/instanceTypes [get]
func EC2InstanceTypes(c *fiber.Ctx) error {
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
	if instanceTypes := amazon.GetEC2InstanceTypes(); instanceTypes != nil {
		return c.Status(fiber.StatusOK).JSON(instanceTypes)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": true,
		"msg":   "Can't find InstanceTypes list in region",
	})
}

func EC2Images(c *fiber.Ctx) error {
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
	if ami := amazon.GetEC2AMI(); ami != nil {
		return c.Status(fiber.StatusOK).JSON(ami)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": true,
		"msg":   "Can't find AMI list in region",
	})
}
