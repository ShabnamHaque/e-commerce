package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const uri = "mongodb://localhost:27017"

func DBSet() *mongo.Client {
	//sets the connection between app and DB.

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	//	err = client.Connect(ctx)
	//err = mongo.Connect(ctx)
	//	client, err := mongo.Connect(ctc.TODO())

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb :(")
		return nil
	}
	fmt.Println("database connected  - Success :)")
	return client

}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productCollection
}
