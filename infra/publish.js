#!/usr/bin/env node

var amqp = require('amqplib/callback_api');
const { nanoid } = require('nanoid');

amqp.connect('amqps://sztwqfjl:la9uwaS03--T93hv0JuJsoiUgxxexMhw@rattlesnake.rmq.cloudamqp.com/sztwqfjl', function(error0, connection) {
  if (error0) {
    throw error0;
  }
  connection.createChannel(function(error1, channel) {
    if (error1) {
      throw error1;
    }
    var exchange = 'direct_xch';
    var routing_key = 'key'

    channel.assertExchange(exchange, 'direct', {
      durable: true
    });
    channel.publish(exchange, routing_key, Buffer.from(JSON.stringify({
        url:"https://vibesmediastorage.s3.amazonaws.com/uploads/61d05316d5d1d2000f61f2d0.mp3",
        type: "metadata",
        uuid: nanoid()
    })));
    console.log(" [x] Sent %s: '%s'", routing_key, 'Ho world!');
  });

  setTimeout(function() {
    connection.close();
    process.exit(0)
  }, 500);
});