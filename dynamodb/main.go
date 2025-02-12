package main

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

func main() {
	dbNew := dynamodb.New(dynamodb.Options{})
}
