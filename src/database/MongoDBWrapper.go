package MongoDBWrapper

import (
	"context"
	"log"
	"time"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect_localhost(URI string, user string, password string) *mongo.Client {

	credential := options.Credential{
		Username: user,
		Password: password,
	}
	//check Authentication 
	clientOptions := options.Client().ApplyURI(URI).SetAuth(credential)
	
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil //failed
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return nil //failed
	}
	return client //succeeded
}

func Connect_client(user string, password string, database string) *mongo.Client {

	clientOptions := options.Client().ApplyURI("mongodb+srv://" + user + ":" + password + "@" + database + ".wmy67.mongodb.net/Cluster0?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
    		log.Fatal(err)
		return nil //failed
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil //failed
	}
	return client //succeeded
}
type MongoFields struct { //fields for demo2 collection
	ID int
	FieldStr string
	FieldInt int
	FieldBool bool
}
func Insert_to_collection(client *mongo.Client, database string, col string, path string){
	
	docsPath, _ := filepath.Abs(path)
	byteValues, err := ioutil.ReadFile(docsPath)
	if err != nil{
		log.Fatal(err)
	}
	if col != "DEMO2" {
		fmt.Println("collection doesn't exist")
		return
	}
	var docs []MongoFields
	err = json.Unmarshal(byteValues, &docs)
	if err != nil{
		log.Fatal(err)
	}
	collection := client.Database(database).Collection(col)

	for i := range docs {
		doc := docs[i]
		result, err := collection.InsertOne(context.TODO(), doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("inserted ID: %v\n", result.InsertedID)

	}
}
