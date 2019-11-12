package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/argandas/serial"
)

var (
	atmhz         *serial.SerialPort
	strhex        string
	line          uint8
	strdataMHZ19B string
)

func getco2(ticker time.Ticker, stop chan bool) {
	var (
		strsuhu string
		strco2  string
		x       uint8
		str     string

		mainHEX string
		THEX    string
		HiHEX   string
		LoHEX   string
		Suhu    int64

		err error
	)

	for {
		select {
		case <-ticker.C:
			strhex = ""

			query := []uint8{255, 001, 134, 000, 000, 000, 000, 000, 121}
			atmhz.Write(query)
			//time.Sleep(50 * time.Millisecond)

			for x = 0; x < 10; x++ {
				line, _ = atmhz.Read()
				//fmt.Printf("%x  ", line)

				str = fmt.Sprintf("%x", line)
				strhex += str + ","
			}
			log.Printf("strhex:%s", strhex)

			if (strings.Contains(strhex, "86") && strings.Contains(strhex, ",0,0,0")) && (strings.Index(strhex, "86") < strings.Index(strhex, ",0,0,0")) {
				mainHEX = strhex[strings.Index(strhex, "86")+3 : strings.Index(strhex, ",0,0,0")]
				THEX = mainHEX[strings.LastIndex(mainHEX, ",")+1 : len(mainHEX)]
				HiHEX = mainHEX[0:strings.Index(mainHEX, ",")]
				LoHEX = mainHEX[strings.Index(mainHEX, ",")+1 : strings.LastIndex(mainHEX, ",")]
				Suhu, err = strconv.ParseInt(THEX, 16, 64)
				if err == nil {
					Suhu -= 40
					//log.Printf("Suhu:%v C", Suhu)
					//log.Printf("tipe data Suhu:%T", Suhu)

					strsuhu = strconv.FormatInt(Suhu, 10)
					//log.Printf("strsuhu:%s C", strsuhu)
				} else {
					strsuhu = "error"
				}

				DecHi, err := strconv.ParseInt(HiHEX, 16, 64)
				if err == nil {
					//log.Printf("DecHi:%v", DecHi)
				}

				DecLi, err := strconv.ParseInt(LoHEX, 16, 64)
				if err == nil {
					//log.Printf("DecLi:%v", DecLi)
				}

				ConCO2 := (DecHi * 256) + DecLi
				//log.Printf("ConCO2:%v ppm", ConCO2)
				//log.Printf("tipe data ConCO2:%T", ConCO2)

				strco2 = strconv.FormatInt(ConCO2, 10)

				strdataMHZ19B = strsuhu + " C|" + strco2 + " ppm"
				messagesco2 <- strdataMHZ19B

			} else {
				log.Println("can't find hex")
			}
			//fmt.Println()

		case <-stop:
			log.Println("stop co2")
			return
		}
	}
}
