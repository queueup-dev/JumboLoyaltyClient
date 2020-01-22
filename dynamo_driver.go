package JumboLoyaltyClient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DatabaseDriver interface {
	saveItem(table string, data interface{}) error
	deleteItem(table string, key string, value string) error
	getItem(table string, key string, value string, object interface{}) (interface{}, error)
}

type dynamoDatabase struct {
	db   *dynamodb.DynamoDB
}

func (d dynamoDatabase) init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	d.db = dynamodb.New(sess)
}

func (d dynamoDatabase) saveItem(table string, data interface{}) error {
	av, err := dynamodbattribute.MarshalMap(data)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:       av,
		TableName:  aws.String(table),
	}

	_, err = d.db.PutItem(input)

	return err
}

func (d dynamoDatabase) deleteItem(table string, key string, value string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			key: {
				S: aws.String(value),
			},
		},
		TableName: aws.String(table),
	}

	_, err := d.db.DeleteItem(input)

	return err
}

func (d dynamoDatabase) getItem(table string, key string, value string, object interface{}) (interface{}, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			key: {
				S: aws.String(value),
			},
		},
	}

	result, err := d.db.GetItem(input)

	if err != nil {
		return object, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &object)

	return object, err
}

func NewDynamoDatabase() *dynamoDatabase {
	db := new(dynamoDatabase)
	db.init()

	return db
}