package queries

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oracle/nosql-go-sdk/nosqldb"
	"github.com/oracle/nosql-go-sdk/nosqldb/jsonutil"
	"github.com/oracle/nosql-go-sdk/nosqldb/types"

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
		return user, err
	}
	if getRes.RowExists() {
		data, err := jsonutil.ToObject(getRes.ValueAsJSON())
		log.Println(getRes.ValueAsJSON())
		if err != nil {
			log.Println(err)
			return user, err
		}
		user.Name = fmt.Sprintf("%v", data["name"])
		user.Email = fmt.Sprintf("%v", data["email"])
		user.Password = fmt.Sprintf("%v", data["password"])
		user.RefreshToken = fmt.Sprintf("%v", data["refreshToken"])
		return user, nil

	}
	log.Println("The row doesn't exist.")
	return user, &fiber.Error{
		Code:    404,
		Message: "The data doesn't exist",
	}
}
