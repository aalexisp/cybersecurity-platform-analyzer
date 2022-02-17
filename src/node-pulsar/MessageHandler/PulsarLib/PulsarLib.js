const Pulsar = require('pulsar-client');

function InitClient(URL){
  const client = new Pulsar.Client({
    serviceUrl: URL,
  });
  return client;
}

function CreateConsumer(client,topic,sub_name){
  const consumer = client.subscribe({
    topic: topic,
    subscription: sub_name,
    subscriptionType: 'Shared',
  });
  return consumer
}

async function ReceiveMessage(consumer){
  const msg = await consumer.receive();
  consumer.acknowledge(msg);
  return msg.getData().toString();
}

function CreateProducer(client,topic){
    const producer = client.createProducer({
        topic: topic,
      });
    return producer;
}

function SendMessage(producer,msg){
  producer.send({
    data: Buffer.from(msg),
  });
  console.log(`Sent message: ${msg}`);
}

module.exports = { InitClient, CreateConsumer, CreateProducer, SendMessage, ReceiveMessage};