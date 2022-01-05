#!/usr/bin/env node

var amqp = require('amqplib/callback_api');
const { nanoid } = require('nanoid');

function getRandomArbitrary(min, max) {
	return Math.random() * (max - min) + min;
}

amqp.connect('amqps://sztwqfjl:la9uwaS03--T93hv0JuJsoiUgxxexMhw@rattlesnake.rmq.cloudamqp.com/sztwqfjl', function (error0, connection) {
  if (error0) {
    throw error0;
  }
  connection.createChannel(function (error1, channel) {
    if (error1) {
      throw error1;
    }
    var exchange = 'direct_xch';
    var routing_key = 'key'

    channel.assertExchange(exchange, 'direct', {
      durable: true
    });
    setInterval(() => {

      channel.publish(exchange, routing_key, Buffer.from(JSON.stringify({
        url: `https://verse.mp3quran.net/arabic/sahl_yassin/64/00${2}${parseInt(getRandomArbitrary(100, 283))}.mp3`,
        type: "metadata",
        uuid: nanoid()
      })));
      console.log(" [x] Sent ");
    }, 2000);
  });

  // setTimeout(function () {
  //   connection.close();
  //   process.exit(0)
  // }, 500);
});