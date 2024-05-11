package persistence

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type userTable struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewUserTable(tableName string, dynamoDB *dynamodb.DynamoDB) *userTable {
	if strings.TrimSpace(tableName) == "" {
		panic("table name is missing")
	}

	if dynamoDB == nil {
		panic("missing dynamodb dependency")
	}

	return &userTable{
		tableName: tableName,
		dynamoDB:  dynamoDB,
	}
}

func (t userTable) CreateUserTable(ctx context.Context) (*dynamodb.DescribeTableOutput, error) {
	_, err := t.dynamoDB.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Email"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(t.tableName),
	})
	if err != nil {
		return nil, err
	}

	tableInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(t.tableName),
	}

	description, err := t.dynamoDB.DescribeTable(tableInput)
	if err != nil {
		return nil, err
	}

	return description, nil
}
