package controllers

import (
	"log"
	"strings"
	"time"

	"main/app/services"
	"main/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// @Router /v1/s3/upload [post]
func UploadHandler(c *fiber.Ctx) error {
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
	log.Println(email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	bucketName, err := services.ConfirmBucketName(strings.Replace(email, "@", ".", 1))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	err = services.UploadToS3(email, bucketName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "upload user tf folder success",
	})
}

// @Router /v1/s3/download [get]
func DownloadHandler(c *fiber.Ctx) error {
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

	bucketName := "terraform-canvas-" + strings.Replace(email, "@", ".", 1)

	zipFilePath, err := services.DownloadToZip(email, bucketName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}

	return c.Status(fiber.StatusOK).SendFile(zipFilePath)
}
