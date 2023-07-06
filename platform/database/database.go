package database

import "github.com/oracle/nosql-go-sdk/nosqldb/auth/iam"

func OCINoSQLConnection() {
	client, err := iam.NewSignatureProvider()
	if err != nil {
		panic(err)
	}
	defer client.Close()
}
