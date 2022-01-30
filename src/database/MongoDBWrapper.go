package MongoDBWrapper

import (
	"context"
	"log"
	"time"
	"fmt"
	"path/filepath"
	"io/ioutil"
//	"reflect"
//	"os"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect_localhost(port string, user string, password string) *mongo.Client {

	credential := options.Credential{
		Username: user,
		Password: password,
	}
	//check authentication

	clientOptions := options.Client().ApplyURI("mongodb://localhost:" + port).SetAuth(credential)
	
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

func Connect_cluster(user string, password string, database string, domain string) *mongo.Client{

	clientOptions := options.Client().ApplyURI("mongodb+srv://" + user + ":" + password + "@" + database + domain)
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

type MongoFields struct {

	ID string `json:"id"`
	FieldStr string `json:"Field Str"`
	FieldInt int `json:"Field Int"`
	FieldBool bool `json:"Field Bool"`
}

func Disconnect(client *mongo.Client){

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
type Filters struct{
	Filter string
	Value string
	Op string
}
type Client_Info struct{
	Client *mongo.Client
	Database string
	Col string

}
func Get_documents(info *Client_Info, fil *[]Filters){
	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)
	var query []bson.M
	for i:= 0;  i<len(*fil); i++ {
		query = append(query, bson.M{"$match": bson.M{(*fil)[i].Filter: bson.M{ (*fil)[i].Op : (*fil)[i].Value}}}) 
	}

	cursor, err := collection.Aggregate(context.TODO(), query)
    	if err != nil {
        	log.Fatal(err)
   	 }
    	
	var docsFiltered []bson.M
	if err = cursor.All(context.TODO(), &docsFiltered); err != nil {
    		log.Fatal(err)
	}
	fmt.Println(docsFiltered)
}

func Insert_to_collection(info *Client_Info, path string) bool{

	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)

	docsPath, _ := filepath.Abs(path)
	byteValues, err := ioutil.ReadFile(docsPath)
	if err != nil{
		log.Fatal(err) //file doesn't exist
		return false
	}
	if (*info).Col != "DEMO3" {
		fmt.Println("collection doesn't exist")
		return false
	}
	var docs  []MongoFields
	err = json.Unmarshal(byteValues, &docs)
	if err != nil{
		log.Fatal(err) //structure is bad
		return false
	}

	for i := range docs {
		doc := docs[i]
		result, err := collection.InsertOne(context.TODO(), doc)
		if err != nil {
			log.Fatal(err)
			return false
		}
		fmt.Printf("inserted ID: %v\n", result.InsertedID)
	}
	return true //succeeded to insert
}
//de moment tots els value son string, nose si necesitem altres tipus
func Delete_documents(info *Client_Info, filter string, value string) bool{

	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	res, err := collection.DeleteMany(ctx, bson.M{ filter : value })
	if err != nil {
 		log.Fatal(err)
		return false
	}
	fmt.Printf("%v document(s) removed\n", res.DeletedCount)
	return true //succeeded to delete
}

func Update_documents(info *Client_Info, filter string, filter_value string,update_Param string, update_value string) bool {

	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := collection.UpdateMany(ctx, bson.M{filter: filter_value},bson.D{
        {"$set", bson.D{{update_Param, update_value}}},},)

	if err != nil {
		log.Fatal(err)
		return false //failed
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	return true //update succeeded
}





