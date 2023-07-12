package OCI

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/objectstorage"

	"main/pkg/configs"
)

func defaultValueSet() (objectstorage.ObjectStorageClient, context.Context, string) {
	client := configs.OCIConfigStorage()
	ctx := context.Background()
	namespace := getNamespace(ctx, client)
	return client, ctx, namespace
}

func getNamespace(ctx context.Context, client objectstorage.ObjectStorageClient) string {
	request := objectstorage.GetNamespaceRequest{}
	r, err := client.GetNamespace(ctx, request)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return *r.Value
}

// 버킷생성
func createBucket(bucketName string) {
	client, ctx, namespace := defaultValueSet()
	request := objectstorage.CreateBucketRequest{
		NamespaceName: &namespace,
	}
	compartmentID := os.Getenv("compartmentID")
	request.CompartmentId = &compartmentID
	request.Name = &bucketName
	request.Metadata = make(map[string]string)
	request.PublicAccessType = objectstorage.CreateBucketDetailsPublicAccessTypeNopublicaccess
	_, err := client.CreateBucket(ctx, request)
	if err != nil {
		log.Println(err)
	}
	log.Println("Created bucket ", bucketName)
}

// 객체 업로드
func PutObject(bucketName, objectName string, contentLen int64, content io.ReadCloser, metadata map[string]string) error {
	client, ctx, namespace := defaultValueSet()
	request := objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
		ContentLength: &contentLen,
		PutObjectBody: content,
		OpcMeta:       metadata,
	}
	_, err := client.PutObject(ctx, request)
	return err
}

// 객체 불러오기
func GetObject(bucketName string, objectName string) ([]string, error) {
	client, ctx, namespace := defaultValueSet()
	request := objectstorage.GetObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
	}
	response, err := client.GetObject(ctx, request)
	buf := new(strings.Builder)
	_, err = io.Copy(buf, response.Content)
	return strings.Split(buf.String(), ","), err
}
