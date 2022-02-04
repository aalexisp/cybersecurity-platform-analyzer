package MongoDBWrapper

import (
	"context"
	"log"
	"time"
	"fmt"
	"path/filepath"
	"io/ioutil"
	//"reflect"
//	"os"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/term"
	"syscall"
	)

type Role struct { //to add roles to the user
    Role string `json:"role"`
    Db   string `json:"db"`
}
type Privilege struct {
	Resource struct {
		Db string `json:"db"`
		Collection string `json:"collection"`
	} `json:"resource"`
	Actions []string `json:"actions"`
}
type Roles struct { //to create new roles
	CreateRole string `json:"createRole"`
	//db string `json:"db"`
	Privileges []Privilege `json:"privileges"`
	Roles []Role `json:"roles"`

}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Roles []Role `json:"roles"`
}
type Machine_definitions struct {
    Machine_name       string `json:"machine_name"`
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

func Add_user(client *mongo.Client, path string){

	docsPath, _ := filepath.Abs(path)
	byteValues, err := ioutil.ReadFile(docsPath)
	if err != nil{
		log.Fatal(err) //file doesn't exist
	}
	var Docs []User

	err = json.Unmarshal(byteValues, &Docs)
	if err != nil{
		log.Fatal(err) //structure is bad
	}
	for i := range Docs {
		doc := Docs[i]

		var rol []bson.M //array of roles that the user has
		for j := 0;  j < len(doc.Roles); j++ {
			rol = append(rol, bson.M{ "role": doc.Roles[j].Role, "db" : doc.Roles[j].Db})
		}
		//we create user in admin database -> we use admin database to auth 
		r := client.Database("admin").RunCommand(context.Background(),bson.D{{"createUser", doc.Username},{"pwd", doc.Password}, {"roles", rol }})

		if r.Err() != nil {
			fmt.Printf("Error creating user %s \n", doc.Username)
			panic(r.Err())

		}else{
			fmt.Printf("User %s was created successfully \n", doc.Username)
		}
	}
}

func Add_role(client *mongo.Client, path string) { //path of the json file with the roles, we create the roles in the admin databse

	docsPath, _ := filepath.Abs(path)
	byteValues, err := ioutil.ReadFile(docsPath)
	if err != nil{
		log.Fatal(err) //file doesn't exist
		return
	//	return false
	}
	var Docs []Roles

	err = json.Unmarshal(byteValues, &Docs)
	if err != nil{
		log.Fatal(err) //structure is bad
		return
	//	return false
	}
	for i := range Docs {
		doc := Docs[i]
	
		var priv bson.A
		for j := range doc.Privileges {//recorrer privilegios
		
			var act bson.A
			for p := range doc.Privileges[j].Actions{
				act = append(act, doc.Privileges[j].Actions[p])
			}
			priv = append(priv, bson.M{"resource": bson.M{ "db": doc.Privileges[j].Resource.Db,"collection": doc.Privileges[j].Resource.Collection ,}, "actions": act ,})

		}
		var rol bson.A
		for t := range doc.Roles {
			rol = append(rol, bson.M{"role": doc.Roles[t].Role, "db":doc.Roles[t].Db})
		}
	//	r := client.Database("admin").RunCommand(context.Background(), bson.D{{"createRole", "newRol"},	{"privileges", bson.A{ bson.M{ "resource": bson.M{"db": "Cluster0", "collection": "",},"actions": bson.A{"insert","dbStats","collStats","compact",},},}}, {"roles", bson.A{bson.M{"role": "readWrite", "db": "Cluster0",},},}})
		r := client.Database("admin").RunCommand(context.Background(), bson.D{{"createRole", doc.CreateRole},{"privileges", priv}, {"roles", rol}})
		if r.Err() != nil {
			fmt.Println("error at runCommand")
			panic(r.Err())
			return
		}else{
			fmt.Printf("Role %s was created successfully\n ", doc.CreateRole)
		}

	}	
}

func Connect(portdomain string, user string, database string, localhost bool) *mongo.Client {

	var clientOptions *options.ClientOptions

	fmt.Println("Please enter the password") //get password
	password,err :=  term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil
	}

	if(localhost == true){
		credential := options.Credential{
			Username: user,
			Password: string(password),
		}
		//check authentication

		clientOptions = options.Client().ApplyURI("mongodb://localhost:" + portdomain + "/" + database).SetAuth(credential)
	}else{

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

func Disconnect(client *mongo.Client){

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
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
	var Docs []Machine_definitions
	if (*info).Col == "Machine_definitions" {
		//var Docs []Machine_definitions
	}else{
		fmt.Println("collection doesn't exist")
		return false
	}
	//var docs []Machine_definitions
	//
	err = json.Unmarshal(byteValues, &Docs)
	if err != nil{
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





