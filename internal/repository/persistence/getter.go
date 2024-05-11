package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/erik-sostenes/auth-api/internal/models"
)

type UserGetter interface {
	Get(context.Context, string) (*models.User, error)
}

type dynamoDBUserGetter struct {
	DynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBUserGetter(tableName string, dynamoDB *dynamodb.DynamoDB) UserGetter {
	if strings.TrimSpace(tableName) == "" {
		panic("table name is missing")
	}

	if dynamoDB == nil {
		panic("dynamoDB dependence is missing")
	}

	return &dynamoDBUserGetter{
		DynamoDB:  dynamoDB,
		tableName: tableName,
	}
}

func (u *dynamoDBUserGetter) Get(ctx context.Context, userID string) (*models.User, error) {
	result, err := u.DynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: &u.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: &userID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("%w: user with id '%s' not fount", models.UserNotFound, userID)
	}

	userDTO := &UserDTO{}

	err = dynamodbattribute.UnmarshalMap(result.Item, userDTO)
	if err != nil {
		return nil, err
	}

	return userDTO.ToDomainUser()
}
