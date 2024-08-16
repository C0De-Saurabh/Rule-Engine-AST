package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client
var client *mongo.Client

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
}

// SaveRule saves a rule and its AST to MongoDB using rule_id
func SaveRule(ruleID string, rule []string) error {
	collection := client.Database("rule_engine").Collection("rules")

	var id interface{}
	objectID, err := primitive.ObjectIDFromHex(ruleID)
	if err == nil {
		// ruleID is a valid ObjectID
		id = objectID
	} else {
		// ruleID is not a valid ObjectID, use as a string
		id = ruleID
	}

	document := bson.M{
		"_id":  id,
		"rule": rule,
	}

	_, err = collection.InsertOne(context.TODO(), document)
	if err != nil {
		return fmt.Errorf("failed to save rule: %w", err)
	}

	return nil
}

// RetrieveRule retrieves a rule's AST from MongoDB using ruleID
func RetrieveRule(ruleID string) ([]string, error) {
	collection := client.Database("rule_engine").Collection("rules")

	var id interface{}
	objectID, err := primitive.ObjectIDFromHex(ruleID)
	if err == nil {
		// ruleID is a valid ObjectID
		id = objectID
	} else {
		// ruleID is not a valid ObjectID, use as a string
		id = ruleID
	}

	filter := bson.M{"_id": id}

	var result struct {
		Rule []string `bson:"rule"`
	}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rule: %w", err)
	}

	return result.Rule, nil
}

// GetMongoClient returns a MongoDB client instance.
func GetMongoClient() (*mongo.Client, error) {
	// Define MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}
