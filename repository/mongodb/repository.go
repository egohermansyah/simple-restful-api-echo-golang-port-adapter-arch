package mongodb

import (
	"context"
	"fmt"
	"golang-vscode-setup/config"
	"golang-vscode-setup/util/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const LOG_IDENTIFIER = "REPOSITORY_MONGODB"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func Connect(config config.MongoDbConfig) (*mongo.Database, error) {
	var auth options.Credential
	uri := fmt.Sprintf("%s://%s:%d", "mongodb", config.Host, config.Port)
	auth.AuthSource = config.Name
	auth.Username = config.Username
	auth.Password = config.Password

	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetAuth(auth))
	if err != nil {
		log.Error(err)
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	return client.Database(config.Name), nil
}

func CreateIndex(collection *mongo.Collection, options *options.IndexOptions, data ...[]string) error {
	for _, values := range data {
		document := bsonx.Doc{}
		for _, value := range values {
			document = append(document, bsonx.Elem{Key: value, Value: bsonx.Int32(1)})
		}
		_, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    document,
				Options: options,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
