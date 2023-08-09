package services_test

import (
	"archive/zip"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"main/app/services"
	"main/platform/amazon"

	"github.com/stretchr/testify/assert"
)

func TestS3ConfigUploadToS3(t *testing.T) {
	existBucket := "exist-bucket" + generateRandomString()
	nonExistBucket := "non-exist-bucket" + generateRandomString()

	_, err := services.ConfirmBucketName(existBucket)
	assert.NoError(t, err)

	nonExist, err := services.ConfirmBucketName(nonExistBucket)
	assert.NoError(t, err)

	check, err := amazon.CheckBucketExists(nonExist)
	assert.NoError(t, err)
	assert.True(t, check, "Bucket name does not exist")

	err = amazon.DeleteS3Bucket("terraform-canvas-" + existBucket)
	assert.NoError(t, err)

	err = amazon.DeleteS3Bucket("terraform-canvas-" + nonExistBucket)
	assert.NoError(t, err)
}

func TestS3UploadToS3(t *testing.T) {
	tmpDir := t.TempDir()
	email := "example@gmail.com"
	bucketName := "terraform-canvas-test-" + generateRandomString()
	uploadDir := filepath.Join(tmpDir, "usertf", email)

	err := os.MkdirAll(uploadDir, os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(uploadDir, "main.tf"), []byte("main"), 0o644)
	assert.NoError(t, err)

	err = os.MkdirAll(uploadDir, os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(uploadDir, "variables.tf"), []byte("variable"), 0o644)
	assert.NoError(t, err)

	err = amazon.CreateBucket(bucketName)
	assert.NoError(t, err)

	err = services.UploadToS3(uploadDir, email, bucketName)
	assert.NoError(t, err)

	objects, err := amazon.ListObjects(bucketName)
	assert.NoError(t, err)

	expectedFiles := []string{email + "/main.tf", email + "/variables.tf"}
	for _, obj := range objects {
		assert.Contains(t, expectedFiles, *obj.Key,
			"Object `"+*obj.Key+"` not found in the folder name `"+email+"`")
	}

	err = amazon.DeleteObjects(bucketName, expectedFiles)
	assert.NoError(t, err)

	err = amazon.DeleteS3Bucket(bucketName)
	assert.NoError(t, err)
}

func TestS3DownloadToZip(t *testing.T) {
	tmpDir := t.TempDir()
	email := "example@gmail.com"
	bucketName := "terraform-canvas-test-" + generateRandomString()
	uploadDir := filepath.Join(tmpDir, "usertf", email)
	downloadDir := filepath.Join(tmpDir, "terraform-canvas")

	err := os.MkdirAll(uploadDir, os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(uploadDir, "main.tf"), []byte("main"), 0o644)
	assert.NoError(t, err)

	err = os.MkdirAll(uploadDir, os.ModePerm)
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(uploadDir, "variables.tf"), []byte("variable"), 0o644)
	assert.NoError(t, err)

	err = amazon.CreateBucket(bucketName)
	assert.NoError(t, err)

	err = services.UploadToS3(uploadDir, email, bucketName)
	assert.NoError(t, err)

	zipFilePath, err := services.DownloadToZip(downloadDir, bucketName)
	assert.NoError(t, err)

	zipFiles, err := zip.OpenReader(zipFilePath)
	assert.NoError(t, err, "Failed to open the zip file")
	defer zipFiles.Close()

	expectedFiles := []string{"main.tf", "variables.tf"}
	for _, file := range zipFiles.File {
		assert.Contains(t, expectedFiles, file.Name,
			"Object `"+file.Name+"` not found in the folder name `"+email+"`")
	}

	err = amazon.DeleteObjects(bucketName, []string{email + "/main.tf", email + "/variables.tf"})
	assert.NoError(t, err)

	err = amazon.DeleteS3Bucket(bucketName)
	assert.NoError(t, err)
}

func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())

	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
