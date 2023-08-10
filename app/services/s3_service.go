package services

import (
	"archive/zip"
	"io"
	"io/ioutil"
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

func UploadToS3(uploadDir string, email string, bucketName string) error {
	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
	   return err
	}
	for _, file := range files {
	   if !file.IsDir() {
		  f, err := os.Open(filepath.Join(uploadDir, file.Name()))
		  key := email + "/" + file.Name()
		  err = amazon.UploadToS3(bucketName, key, f)
		  if err != nil {
			 return err
		  }
	   }
	}
 
	return nil
 }

func DownloadToZip(downloadDir string, bucketName string) (string, error) {
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return "", err
	}

	zipFilePath := downloadDir + ".zip"
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	contents, err := amazon.ListObjects(bucketName)
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
