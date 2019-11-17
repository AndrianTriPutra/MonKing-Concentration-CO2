package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	msgco2 = make(chan string)

	printout string

	strDevice  string
	strDate    string
	strJam     string
	strHeading string
	strLat     string
	strLong    string
	strSpeed   string
	strTemp    string
	strCo2     string
)

func generate(stop chan bool) {
	var (
		gendata string
		strTime string
	)

	for {
		select {
		case gendata = <-msgco2:
			//log.Printf("gendata:%s", gendata)
			strDevice = gendata[strings.Index(gendata, "Device")+10 : strings.Index(gendata, "Time")-6]
			strTime = gendata[strings.Index(gendata, "Time")+8 : strings.Index(gendata, "Heading")-6]
			strHeading = gendata[strings.Index(gendata, "Heading")+11 : strings.Index(gendata, "Latitude")-6]
			strLat = gendata[strings.Index(gendata, "Latitude")+12 : strings.Index(gendata, "Longitude")-6]
			strLong = gendata[strings.Index(gendata, "Longitude")+13 : strings.Index(gendata, "Speed")-6]
			strSpeed = gendata[strings.Index(gendata, "Speed")+9 : strings.Index(gendata, "Temperature")-6]
			strTemp = gendata[strings.Index(gendata, "Temperature")+15 : strings.Index(gendata, "CO2")-6]
			strCo2 = gendata[strings.Index(gendata, "CO2")+7 : strings.Index(gendata, "}")-6]

			strHeading = strHeading[0 : strings.Index(strHeading, "degre")-1]
			strSpeed = strSpeed[0 : strings.Index(strSpeed, "kmh")-1]
			strTemp = strTemp[0 : strings.Index(strTemp, "C")-1]

			strJam = strTime[0:strings.Index(strTime, " ")]
			strDate = strTime[strings.Index(strTime, " ")+1 : len(strTime)]
			strJam = strTime[0:strings.Index(strTime, " ")]

			switch strDevice {
			case "gondril":
				printout = " | " + strDevice + " | " + strDate + " | " + strJam + " | " + strHeading + " | " + strLat + " | " + strLong + " | " +
					strSpeed + " | " + strTemp + " | " + strCo2 + " | "
				//log.Printf("printout:%s", printout)
				sqlInsert()

			default:
				break
			}

			//fmt.Println()

		case <-stop:
			log.Println("stop generate")
			return

		}

	}
}

func sqlInsert() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db-monco2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	tabel := "insert into "
	tabel += strDevice
	tabel += " values (?, ?, ?, ?,?, ?, ?, ?)"

	_, err = db.Exec(tabel, strDate, strJam, strHeading, strLat, strLong, strSpeed, strTemp, strCo2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

/*
{
  "Device": "gondril",
  "Time": "10:24:16 01/10/2019",
  "Heading": "0.00 degre",
  "Latitude": "-6.402677",
  "Longitude": "106.811557",
  "Speed": "0.03 kmh",
  "Temperature": "32 C",
  "CO2": "866 ppm"
 }

*/
