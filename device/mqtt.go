package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	messagestracker = make(chan string)
	messagesco2     = make(chan string)
)

type paket struct {
	Device      string
	Time        string
	Heading     string
	Latitude    string
	Longitude   string
	Speed       string
	Temperature string
	CO2         string
}

var (
	opts   = MQTT.NewClientOptions() //for mqtt client options
	client = MQTT.NewClient(opts)    //for mqtt client
)

var (
	loc                     *time.Location
	timeutcplus7            time.Time //var for utc+7
	hhmmss7, yymmdd7, waktu string    //variabel time +7

	strtracker, strco2 string
	strpaket           string
)

func mqtt() {
	loc, _ = time.LoadLocation("Asia/Jakarta")

	for {
		strtracker, strco2 = <-messagestracker, <-messagesco2
		//log.Printf("strtracker:%s", strtracker)
		//log.Printf("strco2:%s", strco2)
		getTime()

		a := strtracker[0:strings.Index(strtracker, "#")]
		b := strtracker[strings.Index(strtracker, "#")+1 : strings.Index(strtracker, ",")]
		c := strtracker[strings.Index(strtracker, ",")+1 : strings.Index(strtracker, "*")]
		d := strtracker[strings.Index(strtracker, "*")+1 : len(strtracker)]

		e := strco2[0:strings.Index(strco2, "|")]
		f := strco2[strings.Index(strco2, "|")+1 : len(strco2)]

		//log.Printf("a:%s", a)
		//log.Printf("b:%s", b)
		//log.Printf("c:%s", c)
		//log.Printf("d:%s", d)
		//log.Printf("e:%s", e)
		//log.Printf("f:%s", f)

		bufCD := paket{device, waktu, a, b, c, d, e, f}
		WriteCD, err := json.MarshalIndent(bufCD, " ", " ")
		if err == nil {
			strpaket = string(WriteCD)
			strpaket += "\n" //biar enak diliatnya
			//fmt.Printf("%s", strpaket)
			
			//qos0
			sendpaket := client.Publish("co2/monitoring", 0, false, strpaket)
			if sendpaket.WaitTimeout(3*time.Second) != true {
				log.Println("send failed")
			} else {
				fmt.Printf("%s", strpaket)
				//log.Println("send suc")
			}
			/*
			//qos2
			sendpaket := client.Publish("co2/monitoring", 2, false, strpaket)
			if sendpaket.WaitTimeout(3*time.Second) != true {
				log.Println("send failed")
			}*/
		}
	}
}

func getTime() {
	timeutcplus7 = time.Now().In(loc)

	hhmmss7 = fmt.Sprintf("%02d:%02d:%02d", timeutcplus7.Hour(), timeutcplus7.Minute(), timeutcplus7.Second())
	yymmdd7 = fmt.Sprintf("%02d/%02d/%02d", timeutcplus7.Day(), timeutcplus7.Month(), timeutcplus7.Year())
	waktu = hhmmss7 + " " + yymmdd7
}
