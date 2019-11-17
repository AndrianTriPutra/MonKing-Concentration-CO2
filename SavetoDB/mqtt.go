package main

import (
	"log"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	opts   = MQTT.NewClientOptions() //for mqtt client options
	client = MQTT.NewClient(opts)    //for mqtt client
)

var (
	msgSubs string
)

func subscriber(ticker time.Ticker, stop chan bool) {
	for {
		select {
		case <-ticker.C:
			tokenSubs := client.Subscribe(topicmonco2, 0, fhandler)
			if tokenSubs.WaitTimeout(5*time.Second) && tokenSubs.Error() != nil {
				log.Println("Subs on func subs mqtt Err")
			}

		case <-stop:
			log.Println("stop subscriber")
			return
		}
	}
}

var fhandler MQTT.MessageHandler = func(cli MQTT.Client, msg MQTT.Message) {
	topicfhandler := string(msg.Topic())
	if topicfhandler == topicmonco2 {
		msgSubs = string(msg.Payload())
		if (strings.Contains(msgSubs, "Device") && strings.Contains(msgSubs, "Time")) &&
			(strings.Contains(msgSubs, "Heading") && strings.Contains(msgSubs, "Latitude")) &&
			(strings.Contains(msgSubs, "Longitude") && strings.Contains(msgSubs, "Speed")) &&
			(strings.Contains(msgSubs, "Temperature") && strings.Contains(msgSubs, "CO2")) {

			msgco2 <- msgSubs
		}

	}
}
