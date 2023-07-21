package services

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"main/platform/amazon"
)

func ConfirmBucketName(bucketEmail string) (string, error) {
	bucketName := "terraform-canvas-" + bucketEmail

	exists, err := amazon.CheckBucketExists(bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = amazon.CreateBucket(bucketName)
		if err != nil {
			return "", err
		}
	}
	return bucketName, err
}

func UploadToS3(email string, bucketName string) error {
	uploadDir := filepath.Join("usertf", email)

	err := filepath.Walk(uploadDir, func(path string, info os.FileInfo, errWalk error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			key := email + "/" + info.Name()
			err = amazon.UploadToS3(bucketName, key, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func DownloadToZip(email string, bucketName string) (string, error) {
	tempDir := os.TempDir()
	downloadDir := filepath.Join(tempDir, email)
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return "", err
	}

	zipFilePath := filepath.Join(tempDir, email+".zip")
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	contents, err := amazon.ListObjects("terraform-canvas")
	if err != nil {
		return "", err
	}
	for _, obj := range contents {
		key := *obj.Key
		filename := filepath.Join(downloadDir, filepath.Base(key))

		err = amazon.DownloadFile(bucketName, key, filename)
		if err != nil {
			return "", err
		}

		relPath, err := filepath.Rel(downloadDir, filename)
		fileToZip, err := os.Open(filename)
		zipEntry, err := zipWriter.Create(filepath.ToSlash(relPath))
		if err != nil {
			fileToZip.Close()
			return "", err
		}

		_, err = io.Copy(zipEntry, fileToZip)
		if err != nil {
			fileToZip.Close()
			return "", err
		}

		fileToZip.Close()
	}

	zipWriter.Close()
	return zipFilePath, err
}