package persistence

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/erik-sostenes/auth-api/internal/models"
)

type UserSaver interface {
	Save(context.Context, *models.User) error
}

type dynamoDBUserSaver struct {
	DynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBUserSaver(tableName string, dynamoDB *dynamodb.DynamoDB) UserSaver {
	if strings.TrimSpace(tableName) == "" {
		panic("table name is missing")
	}

	if dynamoDB == nil {
		panic("dynamoDB dependence is missing")
	}

	return &dynamoDBUserSaver{
		DynamoDB:  dynamoDB,
		tableName: tableName,
	}
}

func (u *dynamoDBUserSaver) Save(ctx context.Context, user *models.User) error {
	userDTO := &UserDTO{
		Id:            user.ID(),
		Name:          user.Name(),
		Email:         user.Email(),
		Picture:       user.Picture(),
		VerifiedEmail: user.VerifiedEmail(),
	}

	userItem, err := dynamodbattribute.MarshalMap(userDTO)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      userItem,
		TableName: &u.tableName,
	}

	_, err = u.DynamoDB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
