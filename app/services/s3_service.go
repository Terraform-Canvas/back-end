package services

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"main/platform/amazon"
)

func UploadToS3(email string) error {
	uploadDir := filepath.Join("usertf", email)
	client := amazon.GetS3Client()

	err := filepath.Walk(uploadDir, func(path string, info os.FileInfo, errWalk error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			key := email + "/" + info.Name()
			err = amazon.UploadToS3(client, key, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func DownloadToZip(email string) (string, error) {
	client := amazon.GetS3Client()

	resp, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(amazon.GetAWSConfig().BucketName),
		Prefix: aws.String(email + "/"),
	})
	if err != nil {
		return "", err
	}

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

	for _, obj := range resp.Contents {
		key := *obj.Key
		filename := filepath.Join(downloadDir, filepath.Base(key))

		err = amazon.DownloadFile(client, key, filename)
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
