package MachineCreatorWrapper

import (
	"io/ioutil"
	"log"
)


type Machine_definitions struct {
	Machine_name 	  string `json:"Machine_name"`
	Username          string `json:"Username"`
	Password          string `json:"Password"`
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

func ReadMachineBytes(fileName string)  []byte {
	
	byteValue, err := ioutil.ReadFile(fileName)
	if err != nil {
        log.Fatalf("Could not find the file: %v", err)
    }

	return byteValue
}

func WriteMachineToJson(message []byte, filename string) {
	
	err := ioutil.WriteFile(filename, message, 0644)
	if err != nil {
        log.Fatalf("Error creating the file: %v", err)
    }

}
