package MessageHandler

import (
	"db/MessageHandler/PulsarLib"
)

func SendMessage(topic string, message []byte) {
	client := PulsarLib.InitClient("pulsar://localhost:6650")
	defer (*client).Close()
	producer := PulsarLib.CreateProducer(client, topic)
	defer (*producer).Close()
	PulsarLib.SendMessage(producer, message)
}

func RecieveMessage(topic string, sub_name string) []byte {
	client := PulsarLib.InitClient("pulsar://localhost:6650")
	defer (*client).Close()
	consumer := PulsarLib.CreateConsumer(client, topic, sub_name)
	msg := PulsarLib.ReceiveMessage(consumer)
	PulsarLib.DestroyConsumer(consumer)
	return msg
}
