package configs

import (
	"log"
	"os"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

// OCIConfigStorage func for configuration OCI bucket.
func OCIConfigStorage() objectstorage.ObjectStorageClient {
	// Define server settings.
	privateKey, err := os.ReadFile(os.Getenv("privateKeyFile"))
	provider := common.NewRawConfigurationProvider(os.Getenv("tenancyID"),
		os.Getenv("userID"), os.Getenv("region"),
		os.Getenv("fingerprint"), string(privateKey), nil)
	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		log.Println(err)
		return client
	}
	return client
}
