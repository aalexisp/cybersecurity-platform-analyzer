package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar/pulsar-function-go/pf"
)

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

func PublishFunc(ctx context.Context, in []byte) error {
	fctx, ok := pf.FromContext(ctx)
	if !ok {
		return errors.New("get Go Functions Context error")
	}
	var Machines []Machine_definitions
	err := json.Unmarshal(in, &Machines)
	if err != nil {
		log.Fatal(err) //structure is bad
		return nil
	}
	for i := range Machines {
		Machine := &Machines[i]
		Machine.Hello = "bye"
	}
	message, err := json.Marshal(&Machines)
	if err != nil {
		log.Fatal(err) //structure is bad
		return nil
	}
	producer := fctx.NewOutputMessage("test_output")
	msgID, err := producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: message,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("The output message ID is: %+v", msgID)
	return nil
}

func main() {
	pf.Start(PublishFunc)
}
