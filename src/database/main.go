package main

import (
	"db/MessageHandler"
	"db/MongoDBWrapper"
	"fmt"
)

func main() {
	client_localhost := MongoDBWrapper.Connect("27017", "myTester", "admin")
	info := MongoDBWrapper.Client_Info{
		Client:   client_localhost,
		Database: "admin",
		Col:      "Machine_definitions",
	}
	for {
		msg := MessageHandler.RecieveMessage("New_Machines", "test")
		r := MongoDBWrapper.Insert_to_collection(&info, msg)
		if r {
			fmt.Println("Json guardat correctament!")
		}
	}
	MongoDBWrapper.Disconnect(client_localhost)
}
