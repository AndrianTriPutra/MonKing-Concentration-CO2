package main

import (
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	ipaddress   = ""
	port        = ""
	username    = ""
	password    = ""
	topicmonco2 = "co2/monitoring"
)

func setup() {
	dns := "tcp://" + ipaddress + ":" + port
	clientid := "monco2"
	opts = MQTT.NewClientOptions()
	opts.AddBroker(dns).SetClientID(clientid)
	opts.SetUsername(username).SetPassword(password)
	opts.SetMaxReconnectInterval(30 * time.Second).SetDefaultPublishHandler(fhandler)
	client = MQTT.NewClient(opts)
	tokenConn := client.Connect()

	if tokenConn.WaitTimeout(20*time.Second) && tokenConn.Error() != nil {
		log.Println("MQTT Not Connected Start")
	} else {
		log.Println("MQTT Connected Start")
	}
}
