package amazon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"main/pkg/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// 버킷생성
func CreateBucket(bucketName string) error {
	client := configs.GetS3Client()
	if configs.GetAWSConfig().Region != "us-east-1" {
		_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(configs.GetAWSConfig().Region),
			},
		})
		if err != nil {
			return err
		}
	} else {
		_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 버킷 존재여부 확인
func CheckBucketExists(bucketName string) (bool, error) {
	client := configs.GetS3Client()

	_, err := client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		var notFoundErr *types.NotFound
		if errors.As(err, &notFoundErr) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	return true, nil
}

// 버킷 내의 객체 list
func ListObjects(bucketName string) ([]types.Object, error) {
	client := configs.GetS3Client()
	result, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	return result.Contents, err
}

// 업로드
func UploadToS3(bucketName string, key string, file *os.File) error {
	client := configs.GetS3Client()
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	}

	_, err := client.PutObject(context.Background(), params)

	return err
}

// 다운로드
func DownloadFile(bucketName string, key string, destination string) error {
	client := configs.GetS3Client()
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	objResp, err := client.GetObject(context.TODO(), params)
	if err != nil {
		return err
	}
	defer objResp.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, objResp.Body)
	if err != nil {
		return err
	}

	return nil
}

// 버킷 삭제
func DeleteS3Bucket(bucketName string) error {
	client := configs.GetS3Client()
	_, err := client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

// 버킷 내 객체 삭제
func DeleteObjects(bucketName string, objectKeys []string) error {
	client := configs.GetS3Client()
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	return err
}
