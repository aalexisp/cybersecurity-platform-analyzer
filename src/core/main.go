package main

import (
    //"log"
    //"time"
	//"fmt"
	//"context"
    //"github.com/apache/pulsar-client-go/pulsar"
	"cyber-sec/core/MachineCreatorWrapper"
	"cyber-sec/core/PulsarLib"
)

func main() {

	client := PulsarLib.InitClient("pulsar://localhost:6650")
	defer (*client).Close()

	machine := MachineCreatorWrapper.ReadMachineBytes("machine.json")

	producer := PulsarLib.CreateProducer(client, "New_Machines")
	defer (*producer).Close()

	PulsarLib.SendMessage(producer,machine)

}
	

