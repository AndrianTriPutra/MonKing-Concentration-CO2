package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	nmea "github.com/adrianmo/go-nmea"
	"github.com/argandas/serial"
)

var (
	atgps            *serial.SerialPort
	strNumSatellites string
)

func gps(ticker time.Ticker, stop chan bool) {
	var (
		rawline string
		y       uint8
		GNRMC   string
		GNGGA   string
	)
	for {
		select {
		case <-ticker.C:
			GNRMC = ""
			GNGGA = ""

			for y = 0; y < 13; y++ {
				rawline, _ = atgps.ReadLine()

				if strings.Contains(rawline, "$GNRMC") {
					GNRMC = rawline
				} else if strings.Contains(rawline, "$GNGGA") {
					GNGGA = rawline
				}
				//log.Printf("rawline:%s", rawline)
			}
			//log.Printf("GNRMC:%s", GNRMC)
			//log.Printf("GNGGA:%s", GNGGA)

			if len(GNGGA) > 10 {
				m, err := nmea.Parse(GNGGA)
				if err != nil {
					log.Println("err GNGGA")
				} else {
					s := m.(nmea.GNGGA)
					strNumSatellites = strconv.FormatInt(s.NumSatellites, 10)
					//log.Println("NumSatellites: ", strNumSatellites)
				}
			}

			if len(GNRMC) > 10 {
				m, err := nmea.Parse(GNRMC)
				if err != nil {
					log.Println("err GNRMC")
				} else {
					s := m.(nmea.GNRMC)

					//Validity := s.Validity
					speed := s.Speed
					lat := s.Latitude
					long := s.Longitude
					Course := s.Course

					//log.Printf("tipe data speed:%T", speed)
					//log.Printf("tipe data lat:%T", lat)
					//log.Printf("tipe data long:%T", long)
					//log.Printf("tipe data Course:%T", Course)

					if s.Validity == "A" {
						speed *= 1.28
						//log.Printf("Course:%v", Course)
						//log.Printf("speed:%v", speed)
						//log.Printf("lat:%v", lat)
						//log.Printf("long:%v", long)

						strspeed := fmt.Sprintf("%.2f", speed)
						strlat := fmt.Sprintf("%.6f", lat)
						strlong := fmt.Sprintf("%.6f", long)
						strcorse := fmt.Sprintf("%.2f", Course)

						strspeed += " kmh"
						strcorse += " degre"
						//log.Printf("strspeed:%s", strspeed)
						//log.Printf("strlat:%s", strlat)
						//log.Printf("strlong:%s", strlong)
						//log.Printf("strcorse:%s", strcorse)

						datatracker := strcorse + "#" + strlat + "," + strlong + "*" + strspeed
						//log.Printf("datatracker:%s", datatracker)

						messagestracker <- datatracker
					}

				}
			}
			//fmt.Println("======================")
		case <-stop:
			log.Println("stop gps")
			return
		}
	}
}
