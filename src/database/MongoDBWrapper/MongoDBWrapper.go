package MongoDBWrapper

import (
	"context"
	"fmt"
	"log"
	"time"

	//"reflect"
	//	"os"
	"encoding/json"
	"syscall"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/term"
)

type Role struct { //to add roles to the user
	Role string `json:"role"`
	Db   string `json:"db"`
}

type Machine_definitions struct {
	Machine_name      string `json:"machine_name"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	HostIP            string `json:"Host-ip"`
	BoxSpecifications struct {
		CPU      string `json:"Cpu"`
		RAM      string `json:"Ram"`
		Provider string `json:"Provider"`
		TTL      string `json:"Ttl"`
		BoxIP    string `json:"Box-ip"`
		Gui      bool   `json:"Gui"`
	} `json:"Box-specifications"`
	Hello string `json:"Hello"`
}

type Filters struct {
	Filter string
	Value  string
	Op     string
}
type Client_Info struct {
	Client   *mongo.Client
	Database string
	Col      string
}

func Connect(portdomain string, user string, database string) *mongo.Client { //use empty string as database for localhost connection

	var clientOptions *options.ClientOptions

	fmt.Println("Please enter the password") //get password
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil
	}

	if database != "" {
		credential := options.Credential{
			Username: user,
			Password: string(password),
		}
		//check authentication

		clientOptions = options.Client().ApplyURI("mongodb://localhost:" + portdomain + "/" + database).SetAuth(credential)
	} else {

		clientOptions = options.Client().ApplyURI("mongodb+srv://" + user + ":" + string(password) + "@" + database + portdomain)

	}
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

func Disconnect(client *mongo.Client) {

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
func Get_documents(info *Client_Info, fil *[]Filters) {
	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)
	var query []bson.M
	for i := 0; i < len(*fil); i++ {
		query = append(query, bson.M{"$match": bson.M{(*fil)[i].Filter: bson.M{(*fil)[i].Op: (*fil)[i].Value}}})
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

func Insert_to_collection(info *Client_Info, payload []byte) bool {

	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)

	/* 	docsPath, _ := filepath.Abs(path)
	   	byteValues, err := ioutil.ReadFile(docsPath)
	   	if err != nil {
	   		log.Fatal(err) //file doesn't exist
	   		return false
	   	} */
	var Docs []Machine_definitions
	if (*info).Col == "Machine_definitions" {
		//var Docs []Machine_definitions
	} else {
		fmt.Println("collection doesn't exist")
		return false
	}
	//var docs []Machine_definitions
	//
	err := json.Unmarshal(payload, &Docs)
	if err != nil {
		log.Fatal(err) //structure is bad
		return false
	}

	for i := range Docs {
		doc := Docs[i]
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
func Delete_documents(info *Client_Info, filter string, value string) bool {

	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, bson.M{filter: value})
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Printf("%v document(s) removed\n", res.DeletedCount)
	return true //succeeded to delete
}

func Update_documents(info *Client_Info, filter string, filter_value string, update_Param string, update_value string) bool {

	collection := (*info).Client.Database((*info).Database).Collection((*info).Col)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := collection.UpdateMany(ctx, bson.M{filter: filter_value}, bson.D{
		{"$set", bson.D{{update_Param, update_value}}}})

	if err != nil {
		log.Fatal(err)
		return false //failed
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	return true //update succeeded
}
