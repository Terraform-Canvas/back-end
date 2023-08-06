package queries

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/oracle/oci-go-sdk/nosql"

	"main/app/models"
)

type UserQueries struct {
	nosql.NosqlClient
}

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	key := fmt.Sprintf("email:%v", email)
	tableName := "userTable"
	compartmentId := os.Getenv("compartmentID")
	user := models.User{}
	getReq := nosql.GetRowRequest{
		TableNameOrId: &tableName,
		Key:           []string{key},
		CompartmentId: &compartmentId,
	}
	getRes, err := q.GetRow(context.TODO(), getReq)
	if err != nil {
		log.Println(err)
		return user, err
	}
	data := getRes.Value
	user.Name = fmt.Sprintf("%v", data["name"])
	user.Email = fmt.Sprintf("%v", data["email"])
	user.Password = fmt.Sprintf("%v", data["password"])
	user.RefreshToken = fmt.Sprintf("%v", data["refreshToken"])
	return user, nil
}

func (q *UserQueries) UpdateUser(user models.User) error {
	mapValues := map[string]interface{}{
		"email":        user.Email,
		"password":     user.Password,
		"name":         user.Name,
		"refreshToken": user.RefreshToken,
	}
	compartmentId := os.Getenv("compartmentID")
	tableName := "userTable"
	putReq := nosql.UpdateRowRequest{
		TableNameOrId: &tableName,
		UpdateRowDetails: nosql.UpdateRowDetails{
			Value:         mapValues,
			CompartmentId: &compartmentId,
		},
	}

	_, err := q.UpdateRow(context.TODO(), putReq)
	if err != nil {
		return err
	}
	return nil
}
