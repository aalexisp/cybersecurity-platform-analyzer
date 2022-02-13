package main

import (
	"context"
	"db/MongoDBWrapper"
	"db/PulsarLib"
	"fmt"
	"log"
)

func main() {
	client := PulsarLib.InitClient("pulsar://localhost:6650")
	defer (*client).Close()

	consumer := PulsarLib.CreateConsumer(client, "New_Machines", "machine_reader")

	msg, err := (*consumer).Receive(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))

	(*consumer).Ack(msg)

	if err := (*consumer).Unsubscribe(); err != nil {
		log.Fatal(err)
	}

	//MachineCreatorWrapper.WriteMachineToJson(msg.Payload(),"test_done.json") */

	client_localhost := MongoDBWrapper.Connect("27017", "myTester", "admin")
	info := MongoDBWrapper.Client_Info{
		Client:   client_localhost,
		Database: "admin",
		Col:      "Machine_definitions",
	}
	r := MongoDBWrapper.Insert_to_collection(&info, msg.Payload())
	if r {
		fmt.Println("Json guardat correctament!")
	}
	MongoDBWrapper.Disconnect(client_localhost)

}
