package configs

import (
	"log"
	"os"

	"github.com/oracle/oci-go-sdk/nosql"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

// OCIConfigStorage func for configuration OCI bucket.
func OCIConfigStorage() objectstorage.ObjectStorageClient {
	// Define server settings.
	provider := common.NewRawConfigurationProvider(os.Getenv("tenancyID"),
		os.Getenv("userID"), os.Getenv("region"),
		os.Getenv("fingerprint"), os.Getenv("privateKey"), nil)
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		log.Println(err)
		return client
	}
	return client
}

func OCIConfigNoSQL() (nosql.NosqlClient, error) {
	provider := common.NewRawConfigurationProvider(os.Getenv("tenancyID"),
		os.Getenv("userID"), os.Getenv("region"),
		os.Getenv("fingerprint"), os.Getenv("privateKey"), nil)
	client, err := nosql.NewNosqlClientWithConfigurationProvider(provider)
	if err != nil {
		return client, err
	}
	return client, nil
}
