package database

import (
	"main/pkg/configs"

	"main/app/queries"
)

type Queries struct {
	*queries.UserQueries
}

func OCINoSQLConnection() (*Queries, error) {
	client, err := configs.OCIConfigNoSQL()
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{client},
	}, nil
}
