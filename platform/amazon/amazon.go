package amazon

import (
	"context"
	"log"
	"os"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func LoadEnv(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file:", err)
	}
}

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region string
	BucketName string
}

func GetAWSConfig() AWSConfig {
	return AWSConfig{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region: os.Getenv("AWS_REGION"),
		BucketName: os.Getenv("AWS_BUCKET_NAME"),
	}
}

func GetS3Client() (*s3.Client){
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(GetAWSConfig().AccessKey, GetAWSConfig().SecretKey, "")),
		config.WithRegion(GetAWSConfig().Region),
	)
	if err != nil {
		return nil
	}

	return s3.NewFromConfig(cfg)
}

func UploadToS3(client *s3.Client, key string, file *os.File) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(GetAWSConfig().BucketName),
		Key:    aws.String(key),
		Body:   file,
	}

	_, err := client.PutObject(context.Background(), params)
	return err
}

func DownloadFile(client *s3.Client, key string, destination string) error {
	params := &s3.GetObjectInput{
		Bucket: aws.String(GetAWSConfig().BucketName),
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