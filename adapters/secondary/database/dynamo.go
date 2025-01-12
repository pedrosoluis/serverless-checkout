package databaseadapterdynamodb

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/lucasrosa/serverless-checkout/businesslogic/cart"
)

type checkoutRepository struct{}

// NewDynamoCheckoutRepository instantiates the repository for this adapter
func NewDynamoCheckoutRepository() cart.ProcessSecondaryPort {
	return &checkoutRepository{}
}

// PersistedOrder represents the model for inserting the Order into the database
type PersistedOrder struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	ProductID int     `json:"productid"`
}

func (r *checkoutRepository) Save(order *cart.Order) error {
	fmt.Println("saving order", order)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := dynamodb.New(sess)

	persistedOrder := PersistedOrder{
		ID:        order.ID,
		Email:     order.Email,
		Amount:    order.Amount,
		Currency:  order.Currency,
		ProductID: order.ProductID,
	}
	fmt.Println("Persisting order:", persistedOrder)

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(persistedOrder)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Error while sending message to sqs", err)
	} else {
		fmt.Println("Success while sending message to sqs")
	}

	return err
}
