package UsersMongodb

import (
	"context"
	"log"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
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

func Add_role(client *mongo.Client, path string, database string) { //path of the json file with the roles

	docsPath, _ := filepath.Abs(path)
	byteValues, err := ioutil.ReadFile(docsPath)
	if err != nil{
		log.Fatal(err) //file doesn't exist
		return
	}
	var Docs []Roles

	err = json.Unmarshal(byteValues, &Docs)
	if err != nil{
		log.Fatal(err) //structure is bad
		return
	}
	for i := range Docs {
		doc := Docs[i]

		var priv bson.A
		for j := range doc.Privileges { //traverse all privileges

			var act bson.A
			for p := range doc.Privileges[j].Actions{
				act = append(act, doc.Privileges[j].Actions[p])
			}
			priv = append(priv, bson.M{"resource": bson.M{ "db": doc.Privileges[j].Resource.Db,"collection": doc.Privileges[j].Resource.Collection ,}, "actions": act ,})

		}
		if len(doc.Privileges) == 0{
			priv = bson.A{} //no privileges
		}
		var rol bson.A
		for t := range doc.Roles {
			rol = append(rol, bson.M{"role": doc.Roles[t].Role, "db":doc.Roles[t].Db})
		}
		if len(doc.Roles) == 0{
			rol = bson.A{} //no inherited rules
		}
		//r := client.Database("admin").RunCommand(context.Background(), bson.D{{"createRole", "newRol"},	{"privileges", bson.A{ bson.M{ "resource": bson.M{"db": "Cluster0", "collection": "",},"actions": bson.A{"insert","dbStats","collStats","compact",},},}}, {"roles", bson.A{bson.M{"role": "readWrite", "db": "Cluster0",},},}})
		r := client.Database(database).RunCommand(context.Background(), bson.D{{"createRole", doc.CreateRole},{"privileges", priv}, {"roles", rol, }})
		if r.Err() != nil {
			fmt.Println("error at runCommand")
			panic(r.Err())
			return
		}else{
			fmt.Printf("Role %s was created successfully\n ", doc.CreateRole)
		}
	}

}

