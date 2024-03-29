package main

import (
	"log"
	"time"

	"github.com/argandas/serial"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	ipaddress = ""
	port      = ""
	username  = ""
	password  = ""
	device    = "gondril"

	dns      string
	clientid string
)

func setup() {
	//open port mhz
	atmhz = serial.New()
	err := atmhz.Open("/dev/ttyUSB0", 9600, 1*time.Second)
	if err != nil {
		log.Println("PORT MHZ BUSY")
	} else {
		log.Println("SUCCESS OPEN PORT MHZ")
	}

	//open port gps
	atgps = serial.New()
	err = atgps.Open("/dev/ttyACM0", 9600, 1*time.Second)
	if err != nil {
		log.Println("PORT gps BUSY")
	} else {
		log.Println("SUCCESS OPEN PORT gps")
	}

	//mqtt
	dns = "tcp://" + ipaddress + ":" + port
	clientid = "co2/" + device

	opts = MQTT.NewClientOptions()
	opts.AddBroker(dns).SetClientID(clientid)
	opts.SetUsername(username).SetPassword(password)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
        //opts.SetCleanSession(false).SetStore(MQTT.NewFileStore("mqttstorage"))//if you use qos2, make directory mqttstorage
	opts.SetWill(device+"/wills", "good-bye!", 0, false)

	client = MQTT.NewClient(opts)
	tokenConn := client.Connect()

	if tokenConn.WaitTimeout(20*time.Second) && tokenConn.Error() != nil {
		log.Println("MQTT Not Connected")
	} else {
		log.Println("MQTT Connected")
	}

}
