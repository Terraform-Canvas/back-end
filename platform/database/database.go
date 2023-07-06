package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oracle/nosql-go-sdk/nosqldb/common"

	"github.com/oracle/nosql-go-sdk/nosqldb"
	"github.com/oracle/nosql-go-sdk/nosqldb/auth/iam"
	"github.com/oracle/nosql-go-sdk/nosqldb/jsonutil"
	"github.com/oracle/nosql-go-sdk/nosqldb/logger"
	"github.com/oracle/nosql-go-sdk/nosqldb/types"
)

func OCINoSQLConnection() {
	client, err := createClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
}

func createClient() (*nosqldb.Client, error) {
	// provider, err := iam.NewSignatureProviderFromFile("~/.oci/config", "DEFAULT", "",os.Getenv("compartmentID"))
	privateKeyPass := os.Getenv("privateKeyPassphrase")
	provider, err := iam.NewRawSignatureProvider(os.Getenv("tenancyID"),
		os.Getenv("userID"), os.Getenv("region"),
		os.Getenv("fingerprint"), os.Getenv("compartmentID"), os.Getenv("privateKeyFile"), &privateKeyPass)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	cfg := nosqldb.Config{
		Mode:                  "cloud",
		Region:                common.Region(os.Getenv("region")),
		AuthorizationProvider: provider,
	}
	cfg.LoggingConfig = nosqldb.LoggingConfig{
		Logger: logger.New(os.Stdout, logger.Warn, false),
	}
	client, err := nosqldb.NewClient(cfg)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return client, nil
}

func putData(client *nosqldb.Client, tableName string) {
	mapValues := types.ToMapValue("id", 1)
	mapValues.Put("user_id", "test")
	mapValues.Put("password", "testPWD")

	putReq := &nosqldb.PutRequest{
		TableName: tableName,
		Value:     mapValues,
	}
	putRes, err := client.Put(putReq)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println(putRes)
}

func delData(client *nosqldb.Client, tableName string) {
	key := types.ToMapValue("id", 1)
	delReq := &nosqldb.DeleteRequest{
		TableName: tableName,
		Key:       key,
	}
	delRes, err := client.Delete(delReq)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if delRes.Success {
		log.Println("result : " + jsonutil.AsJSON(delRes))
	}
}

func getData(client *nosqldb.Client, tableName string) {
	key := types.ToMapValue("id", 1)
	getReq := &nosqldb.GetRequest{
		TableName: tableName,
		Key:       key,
	}
	getRes, err := client.Get(getReq)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if getRes.RowExists() {
		log.Println(getRes.ValueAsJSON())
	} else {
		log.Println("The row doesn't exist.")
	}
}

func createTable(client *nosqldb.Client, tableName string) {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"id LONG, "+
		"user_id STRING, "+
		"password STRING, "+
		"PRIMARY KEY(id))", tableName)
	tableReq := &nosqldb.TableRequest{
		TableLimits: &nosqldb.TableLimits{
			ReadUnits:  50,
			WriteUnits: 50,
			StorageGB:  2,
		},
		Statement: stmt,
	}
	tableRes, err := client.DoTableRequest(tableReq)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	_, err = tableRes.WaitForCompletion(client, 60*time.Second, time.Second)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func dropTable(client *nosqldb.Client, tableName string) {
	dropReq := &nosqldb.TableRequest{
		Statement: "DROP TABLE IF EXISTS " + tableName,
	}
	tableRes, err := client.DoTableRequestAndWait(dropReq, 60*time.Second, time.Second)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println(tableRes)
}
