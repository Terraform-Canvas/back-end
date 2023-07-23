package controllers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"main/pkg/utils"
	"main/platform/OCI"

	"github.com/gofiber/fiber/v2"

	"main/app/models"
	"main/platform/database"
)

// UserSignIn method to auth user and return access and refresh tokens.
// @Router /v1/login/new [post]
func UserSignIn(c *fiber.Ctx) error {
	// Return status 200 OK.
	signIn := &models.SignIn{}

	if err := c.BodyParser(signIn); err != nil {
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

// @Router /v1/login/refresh [post]
func UserRefresh(c *fiber.Ctx) error {
	now := time.Now().Unix()

	refreshToken := c.Get("X-refresh-token")
	expiresRefreshToken, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	// 시간 확인
	if now < expiresRefreshToken {
		email, err := utils.GetEmailFromToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   err,
			})
		}
		// 토큰 재발급
		tokens, err := utils.GenerateNewTokens(email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "Failed to generate tokens",
			})
		}
		db, err := database.OCINoSQLConnection()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err,
			})
		}
		user, err := db.GetUserByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   err,
			})
		}
		user.RefreshToken = tokens.Refresh
		// db 업데이트
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
		})
	} else {
		// 세션 만료 고지
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, your session was ended earlier",
		})
	}
}

// UserKeySave는 key를 bucket에 업로드하고 env로 해당 키를 등록하는 함수
// @Router /v1/user/key [post]
func UserKeySave(c *fiber.Ctx) error {
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
	userKey := &models.UserKey{}
	if err := c.BodyParser(userKey); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	filePath := fmt.Sprintf("./%s.csv", userKey.Email)
	err = createFile(filePath, userKey.AccessKey, userKey.SecretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	file, err := os.Open(filePath)
	fi, err := file.Stat()
	if err != nil {
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

func createFile(filePath string, accessKey string, secretKey string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// csv writer 생성
	wr := csv.NewWriter(bufio.NewWriter(file))
	// csv 내용 쓰기
	wr.Write([]string{accessKey, secretKey})
	wr.Flush()
	return nil
}
