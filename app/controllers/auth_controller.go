package controllers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"main/pkg/utils"
	"main/platform/OCI"

	"github.com/gofiber/fiber/v2"

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	db, err := database.OCINoSQLConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	user, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	if user.Password != signIn.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given password is not correct",
		})
	}
	tokens, err := utils.GenerateNewTokens(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to generate tokens",
		})
	}

	user.RefreshToken = tokens.Refresh
	err = db.UpdateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to update user refreshToken",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"authInfo": fiber.Map{
			"accessToken":  tokens.Access,
			"refreshToken": tokens.Refresh,
		},
		"name":  user.Name,
		"email": user.Email,
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Router /v1/user/sign/out [post]
func UserSignOut(c *fiber.Ctx) error {
	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}

// UserKeySave는 key를 bucket에 업로드하고 env로 해당 키를 등록하는 함수
// @Router /v1/user/key [post]
func UserKeySave(c *fiber.Ctx) error {
	userKey := &models.UserKey{}
	if err := c.BodyParser(userKey); err != nil {
		log.Println(userKey)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	filePath := fmt.Sprintf("./%s.csv", userKey.Email)
	createFile(filePath, userKey.AccessKey, userKey.SecretKey)
	file, err := os.Open(filePath)
	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	err = OCI.PutObject("canvas-bucket", userKey.Email+".csv", fi.Size(), file, nil)
	defer func() {
		os.Remove(filePath)
	}()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	result, err := OCI.GetObject("canvas-bucket", userKey.Email+".csv")
	result = deleteWithTrim(result)
	// 환경변수에 등록 (추후 변경이 필요한 부분 => 이 키를 사용하는 부분은 버킷에서 불러오는게 효율적인 방법일듯)
	os.Setenv("AWS_ACCESS_KEY", result[0])
	os.Setenv("AWS_SECRET_KEY", result[1])
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "userKeySave success",
	})
}

func deleteWithTrim(result []string) []string {
	for idx, key := range result {
		result[idx] = strings.Trim(key, "\n")
	}
	return result
}

func createFile(filePath string, accessKey string, secretKey string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
	}
	// csv writer 생성
	wr := csv.NewWriter(bufio.NewWriter(file))
	// csv 내용 쓰기
	wr.Write([]string{accessKey, secretKey})
	wr.Flush()
}