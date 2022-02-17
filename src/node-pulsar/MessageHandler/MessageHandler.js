const PulsarLib = require("./PulsarLib/PulsarLib.js");

async function sendMessage(topic, msg){
    // Create a client
    const client = PulsarLib.InitClient("pulsar://localhost:6650");
    // Create a producer
    const producer = await PulsarLib.CreateProducer(client,topic);
    // Send Message
    PulsarLib.SendMessage(producer,msg);
    await producer.flush();
    await producer.close();
    await client.close();
}

async function receiveMessage(topic, sub_name) {
    // Create a client
    const client = PulsarLib.InitClient("pulsar://localhost:6650");
    // Create a consumer
    const consumer = await PulsarLib.CreateConsumer(client,topic,sub_name);
    // Receive messages
    const msg = await PulsarLib.ReceiveMessage(consumer);
    await consumer.close();
    await client.close();
    return msg
}

module.exports = {sendMessage, receiveMessage}