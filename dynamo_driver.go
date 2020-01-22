package JumboLoyaltyClient

import (
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DatabaseDriver interface {
	saveItem(table string, data interface{}) error
	deleteItem(table string, key string, value string) error
	getItem(table string, key string, value string, object interface{}) (interface{}, error)
	listItems(table string, conditions map[string]string, objects interface{}) (interface{}, error)
}

var (
	Dynamo *dynamoDatabase
	database = dynamodb.New(awsSession.Must(awsSession.NewSessionWithOptions(awsSession.Options{
		SharedConfigState: awsSession.SharedConfigEnable,
	})))
)

type dynamoDatabase struct {
	db      *dynamodb.DynamoDB
	Session *awsSession.Session
}

func (d dynamoDatabase) saveItem(table string, data interface{}) error {
	av, err := dynamodbattribute.MarshalMap(data)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
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

// @todo we should add more options then key,val EQ.
func (d dynamoDatabase) listItems(
	table string,
	index string,
	conditions map[string]string,
	objects interface{},
) (interface{}, error) {

	queryConditions := make(map[string]*dynamodb.Condition)
	for key, value := range conditions {
		condition := dynamodb.Condition{
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dynamodb.AttributeValue{
				{
					S: aws.String(value),
				},
			},
		}

		queryConditions[key] = &condition
	}

	queryInput := &dynamodb.QueryInput{
		IndexName: aws.String(index),
		KeyConditions: queryConditions,
		TableName:     aws.String(table),
	}

	result, err := d.db.Query(queryInput)

	if err != nil {
		return nil, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, objects)

	return objects, err
}

func init() {
	Dynamo = new(dynamoDatabase)
	Dynamo.db = database
}