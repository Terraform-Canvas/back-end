package queries

import (
	"fmt"
	"github.com/oracle/nosql-go-sdk/nosqldb"
	"github.com/oracle/nosql-go-sdk/nosqldb/jsonutil"
	"github.com/oracle/nosql-go-sdk/nosqldb/types"
	"log"
	"main/app/models"
)

type UserQueries struct {
	*nosqldb.Client
}

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	key := types.ToMapValue("email", email)
	user := models.User{}
	getReq := &nosqldb.GetRequest{
		TableName: "userTable",
		Key:       key,
	}
	getRes, err := q.Get(getReq)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if getRes.RowExists() {
		data, err := jsonutil.ToObject(getRes.ValueAsJSON())
		log.Println(getRes.ValueAsJSON())
		if err != nil {
			log.Println(err)
			return user, err
		}
		user.Email = fmt.Sprintf("%v", data["email"])
		user.Email = fmt.Sprintf("%v", data["password"])
		return user, nil

	}
	log.Println("The row doesn't exist.")
	return user, fmt.Errorf("The row doesn't exist.")
}
