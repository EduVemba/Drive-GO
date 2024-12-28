package MongoDB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "Uber_Go"
const coolRequest = "Request"
const coolResponse = "Response"

var (
	client      *mongo.Client
	Collections = make(map[string]*mongo.Collection)
)

func Connect() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}
	Collections["requests"] = client.Database(dbName).Collection(coolRequest)
	Collections["responses"] = client.Database(dbName).Collection(coolResponse)

	fmt.Println("Connected to MongoDB!")
}

func InsertOneRequest(content interface{}) error {
	collection := Collections["requests"]
	if collection == nil {
		return fmt.Errorf("collection não inicializada")
	}

	_, err := collection.InsertOne(context.Background(), content)
	if err != nil {
		return fmt.Errorf("error trying to insert: %v", err)
	}

	return nil
}

func InsertOneResponse(content interface{}) error {
	collection := Collections["responses"]
	if collection == nil {
		return fmt.Errorf("collection não inicializada")
	}

	_, err := collection.InsertOne(context.Background(), content)
	if err != nil {
		return fmt.Errorf("error trying to insert: %v", err)
	}

	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Collections[collectionName]
}

func Close() {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Printf("Erro ao desconectar do MongoDB: %v", err)
		} else {
			log.Println("Conexão com MongoDB fechada com sucesso.")
		}
	}
}
