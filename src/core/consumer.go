package main

import(
	"cyber-sec/core/MachineCreatorWrapper"
	"cyber-sec/core/PulsarLib"
	"fmt"
	"log"
	"context"
)

func main (){
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

	MachineCreatorWrapper.WriteMachineToJson(msg.Payload(),"test_done.json")

}