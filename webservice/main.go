package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type datamon struct {
	Tanggal     string
	Jam         string
	Heading     string
	Latitude    string
	Longitude   string
	Speed       string
	Temperature string
	CO2         string
}

type databel struct {
	Line1, Line2, Line3, Line4, Line5, Line6, Line7, Line8, Line9, Line10          [9]string
	Line11, Line12, Line13, Line14, Line15, Line16, Line17, Line18, Line19, Line20 [9]string
	Line21, Line22, Line23, Line24, Line25, Line26, Line27, Line28, Line29, Line30 [9]string
	Line31, Line32, Line33, Line34, Line35, Line36, Line37, Line38, Line39, Line40 [9]string
	Line41, Line42, Line43, Line44, Line45, Line46, Line47, Line48, Line49, Line50 [9]string
	Line51, Line52, Line53, Line54, Line55, Line56, Line57, Line58, Line59, Line60 [9]string
}

var (
	i         = 0
	data      = databel{}
	result    = datamon{}
	strnumber [60]string

	arrLine1, arrLine2, arrLine3, arrLine4, arrLine5, arrLine6, arrLine7, arrLine8, arrLine9, arrLine10          [8]string
	arrLine11, arrLine12, arrLine13, arrLine14, arrLine15, arrLine16, arrLine17, arrLine18, arrLine19, arrLine20 [8]string
	arrLine21, arrLine22, arrLine23, arrLine24, arrLine25, arrLine26, arrLine27, arrLine28, arrLine29, arrLine30 [8]string
	arrLine31, arrLine32, arrLine33, arrLine34, arrLine35, arrLine36, arrLine37, arrLine38, arrLine39, arrLine40 [8]string
	arrLine41, arrLine42, arrLine43, arrLine44, arrLine45, arrLine46, arrLine47, arrLine48, arrLine49, arrLine50 [8]string
	arrLine51, arrLine52, arrLine53, arrLine54, arrLine55, arrLine56, arrLine57, arrLine58, arrLine59, arrLine60 [8]string
)

var (
	tgl   [1000]string
	jam   [1000]string
	lati  [1000]string
	longi [1000]string
	kec   [1000]string
	suhu  [1000]string
	emco2 [1000]int

	max, min           int
	tempmax, tempmin   string
	latmax, longmax    string
	latmin, longmin    string
	datemin, jamin     string
	datemax, jamax     string
	speedmax, speedmin string

	tgl17   [1000]string
	jam17   [1000]string
	lati17  [1000]string
	longi17 [1000]string
	kec17   [1000]string
	suhu17  [1000]string
	emco217 [1000]int

	max17, min17           int
	tempmax17, tempmin17   string
	latmax17, longmax17    string
	latmin17, longmin17    string
	datemin17, jamin17     string
	datemax17, jamax17     string
	speedmax17, speedmin17 string
)

var (
	loc                     *time.Location
	timeutcplus7            time.Time //var for utc+7
	hhmmss7, yymmdd7, waktu string    //variabel time +7
)

func main() {
	fmt.Println("Fimrware Version:29/10/2019")
	fmt.Println("server started at localhost:8080 ")
	fmt.Printf("       look at http://localhost:8080/monstrackco2/home")

	loc, _ = time.LoadLocation("Asia/Jakarta")

	generatemap17()
	generatemap29()

	http.HandleFunc("/monstrackco2/home", routehome)
	http.HandleFunc("/monstrackco2/table", routetable)
	http.HandleFunc("/monstrackco2/map17", routemap17)
	http.HandleFunc("/monstrackco2/map29", routemap29)

	var wg sync.WaitGroup
	stopchan := make(chan bool, 1)

	go func() {
		listen()
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		showtime(*ticker, stopchan)

	}()

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, syscall.SIGINT, syscall.SIGTERM)
	<-cc

	close(stopchan)
	wg.Wait()

}

func listen() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func showtime(ticker time.Ticker, stop chan bool) {
	for {
		select {
		case <-ticker.C:
			sqlQueryRow()

		case <-stop:
			log.Println("stop showtime")
			return
		}
	}
}

func getTime() {
	timeutcplus7 = time.Now().In(loc)

	hhmmss7 = fmt.Sprintf("%02d:%02d:%02d", timeutcplus7.Hour(), timeutcplus7.Minute(), timeutcplus7.Second())
	yymmdd7 = fmt.Sprintf("%02d/%02d/%02d", timeutcplus7.Day(), timeutcplus7.Month(), timeutcplus7.Year())
	waktu = hhmmss7 + " " + yymmdd7
}

func sqlQueryRow() {
	//fmt.Println("i'm query")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db-monco2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	getTime()
	//yymmdd7 = "04/11/2019"
	//log.Printf("yymmdd7:%s", yymmdd7)//yymmdd7:12/11/2019
	//query di sini kurang bagus, oerlu diperbaiki cara query atau manajeman db-nya
	query := "SELECT * FROM `db-monco2`.gondril WHERE Tanggal LIKE "
	query += "'" + yymmdd7 + "%';"

	rows, err := db.Query(query)
	//SELECT * FROM `db-monco2`.gondril ORDER BY Tanggal DESC;
	if err != nil {
		fmt.Println(err.Error())
	} else {
		i = 0
		for rows.Next() {
			err := rows.Scan(&result.Tanggal, &result.Jam, &result.Heading, &result.Latitude, &result.Longitude, &result.Speed, &result.Temperature, &result.CO2)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				bufferResult(i)
				if i < 60 {
					strnumber[i] = strconv.Itoa(i + 1)
				} else {
					break
				}
				//fmt.Printf("| %v | %s | %s | %s | %s | %s | %s | %s | %s |\n", i+1, result.Tanggal, result.Jam, result.Heading, result.Latitude, result.Longitude, result.Speed, result.Temperature, result.CO2)
				i++
			}
		}
	}

}

func routetable(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("table").ParseFiles("view/tabel.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		getline()

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func generatemap17() {
	i := 0
	max17 = 0
	min17 = 2000
	for i = 0; i < 1000; i++ {
		tgl17[i] = ""
		jam17[i] = ""
		lati17[i] = ""
		longi17[i] = ""
		kec17[i] = ""
		suhu17[i] = ""
		emco217[i] = 0
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db-monco2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `db-monco2`.gondril WHERE Tanggal LIKE '17/10/2019%';")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		i = 0
		for rows.Next() {
			err := rows.Scan(&result.Tanggal, &result.Jam, &result.Heading, &result.Latitude, &result.Longitude, &result.Speed, &result.Temperature, &result.CO2)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				//fmt.Printf("| %v | %s | %s | %s | %s | %s | %s | %s | %s |\n", i+1, result.Tanggal, result.Jam, result.Heading, result.Latitude, result.Longitude, result.Speed, result.Temperature, result.CO2)
				buffer := result.CO2[0 : len(result.CO2)-1]
				dataco2, _ := strconv.Atoi(buffer)

				tgl17[i] = result.Tanggal
				jam17[i] = result.Jam
				lati17[i] = result.Latitude
				longi17[i] = result.Longitude
				kec17[i] = result.Speed
				suhu17[i] = result.Temperature
				emco217[i] = dataco2

				if max17 < emco217[i] {
					max17 = emco217[i]
					latmax17 = lati17[i]
					longmax17 = longi17[i]
					datemax17 = tgl17[i]
					jamax17 = jam17[i]
					tempmax17 = suhu17[i]
					speedmax17 = kec17[i]
				}

				if min17 > emco217[i] {
					min17 = emco217[i]
					latmin17 = lati17[i]
					longmin17 = longi17[i]
					datemin17 = tgl17[i]
					jamin17 = jam17[i]
					tempmin17 = suhu17[i]
					speedmin17 = kec17[i]
				}
				i++
			}
		}
	}
}

func generatemap29() {
	i := 0
	max = 0
	min = 2000
	for i = 0; i < 1000; i++ {
		tgl[i] = ""
		jam[i] = ""
		lati[i] = ""
		longi[i] = ""
		kec[i] = ""
		suhu[i] = ""
		emco2[i] = 0
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db-monco2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `db-monco2`.gondril WHERE Tanggal LIKE '29/10/2019%';")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		i = 0
		for rows.Next() {
			err := rows.Scan(&result.Tanggal, &result.Jam, &result.Heading, &result.Latitude, &result.Longitude, &result.Speed, &result.Temperature, &result.CO2)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				//fmt.Printf("| %v | %s | %s | %s | %s | %s | %s | %s | %s |\n", i+1, result.Tanggal, result.Jam, result.Heading, result.Latitude, result.Longitude, result.Speed, result.Temperature, result.CO2)
				buffer := result.CO2[0 : len(result.CO2)-1]
				dataco2, _ := strconv.Atoi(buffer)

				tgl[i] = result.Tanggal
				jam[i] = result.Jam
				lati[i] = result.Latitude
				longi[i] = result.Longitude
				kec[i] = result.Speed
				suhu[i] = result.Temperature
				emco2[i] = dataco2

				if max < emco2[i] {
					max = emco2[i]
					latmax = lati[i]
					longmax = longi[i]
					datemax = tgl[i]
					jamax = jam[i]
					tempmax = suhu[i]
					speedmax = kec[i]
				}

				if min > emco2[i] {
					min = emco2[i]
					latmin = lati[i]
					longmin = longi[i]
					datemin = tgl[i]
					jamin = jam[i]
					tempmin = suhu[i]
					speedmin = kec[i]
				}
				i++
			}
		}
	}

}

func getline() {
	data = databel{
		Line1: [9]string{strnumber[0], arrLine1[0], arrLine1[1], arrLine1[2],
			arrLine1[3], arrLine1[4], arrLine1[5], arrLine1[6], arrLine1[7]},
		Line2: [9]string{strnumber[1], arrLine2[0], arrLine2[1], arrLine2[2],
			arrLine2[3], arrLine2[4], arrLine2[5], arrLine2[6], arrLine2[7]},
		Line3: [9]string{strnumber[2], arrLine3[0], arrLine3[1], arrLine3[2],
			arrLine3[3], arrLine3[4], arrLine3[5], arrLine3[6], arrLine3[7]},
		Line4: [9]string{strnumber[3], arrLine4[0], arrLine4[1], arrLine4[2],
			arrLine4[3], arrLine4[4], arrLine4[5], arrLine4[6], arrLine4[7]},
		Line5: [9]string{strnumber[4], arrLine5[0], arrLine5[1], arrLine5[2],
			arrLine5[3], arrLine5[4], arrLine5[5], arrLine5[6], arrLine5[7]},
		Line6: [9]string{strnumber[5], arrLine6[0], arrLine6[1], arrLine6[2],
			arrLine6[3], arrLine6[4], arrLine6[5], arrLine6[6], arrLine6[7]},
		Line7: [9]string{strnumber[6], arrLine7[0], arrLine7[1], arrLine7[2],
			arrLine7[3], arrLine7[4], arrLine7[5], arrLine7[6], arrLine7[7]},
		Line8: [9]string{strnumber[7], arrLine8[0], arrLine8[1], arrLine8[2],
			arrLine8[3], arrLine8[4], arrLine8[5], arrLine8[6], arrLine8[7]},
		Line9: [9]string{strnumber[8], arrLine9[0], arrLine9[1], arrLine9[2],
			arrLine9[3], arrLine9[4], arrLine9[5], arrLine9[6], arrLine9[7]},
		Line10: [9]string{strnumber[9], arrLine10[0], arrLine10[1], arrLine10[2],
			arrLine10[3], arrLine10[4], arrLine10[5], arrLine10[6], arrLine10[7]},
		Line11: [9]string{strnumber[10], arrLine11[0], arrLine11[1], arrLine11[2],
			arrLine11[3], arrLine11[4], arrLine11[5], arrLine11[6], arrLine11[7]},
		Line12: [9]string{strnumber[11], arrLine12[0], arrLine12[1], arrLine12[2],
			arrLine12[3], arrLine12[4], arrLine12[5], arrLine12[6], arrLine12[7]},
		Line13: [9]string{strnumber[12], arrLine13[0], arrLine13[1], arrLine13[2],
			arrLine13[3], arrLine13[4], arrLine13[5], arrLine13[6], arrLine13[7]},
		Line14: [9]string{strnumber[13], arrLine14[0], arrLine14[1], arrLine14[2],
			arrLine14[3], arrLine14[4], arrLine14[5], arrLine14[6], arrLine14[7]},
		Line15: [9]string{strnumber[14], arrLine15[0], arrLine15[1], arrLine15[2],
			arrLine15[3], arrLine15[4], arrLine15[5], arrLine15[6], arrLine15[7]},
		Line16: [9]string{strnumber[15], arrLine16[0], arrLine16[1], arrLine16[2],
			arrLine16[3], arrLine16[4], arrLine16[5], arrLine16[6], arrLine16[7]},
		Line17: [9]string{strnumber[16], arrLine17[0], arrLine17[1], arrLine17[2],
			arrLine17[3], arrLine17[4], arrLine17[5], arrLine17[6], arrLine17[7]},
		Line18: [9]string{strnumber[17], arrLine18[0], arrLine18[1], arrLine18[2],
			arrLine18[3], arrLine18[4], arrLine18[5], arrLine18[6], arrLine18[7]},
		Line19: [9]string{strnumber[18], arrLine19[0], arrLine19[1], arrLine19[2],
			arrLine19[3], arrLine19[4], arrLine19[5], arrLine19[6], arrLine19[7]},
		Line20: [9]string{strnumber[19], arrLine20[0], arrLine20[1], arrLine20[2],
			arrLine20[3], arrLine20[4], arrLine20[5], arrLine20[6], arrLine20[7]},
		Line21: [9]string{strnumber[20], arrLine21[0], arrLine21[1], arrLine21[2],
			arrLine21[3], arrLine21[4], arrLine21[5], arrLine21[6], arrLine21[7]},
		Line22: [9]string{strnumber[21], arrLine22[0], arrLine22[1], arrLine22[2],
			arrLine22[3], arrLine22[4], arrLine22[5], arrLine22[6], arrLine22[7]},
		Line23: [9]string{strnumber[22], arrLine23[0], arrLine23[1], arrLine23[2],
			arrLine23[3], arrLine23[4], arrLine23[5], arrLine23[6], arrLine23[7]},
		Line24: [9]string{strnumber[23], arrLine24[0], arrLine24[1], arrLine24[2],
			arrLine24[3], arrLine24[4], arrLine24[5], arrLine24[6], arrLine24[7]},
		Line25: [9]string{strnumber[24], arrLine25[0], arrLine25[1], arrLine25[2],
			arrLine25[3], arrLine25[4], arrLine25[5], arrLine25[6], arrLine25[7]},
		Line26: [9]string{strnumber[25], arrLine26[0], arrLine26[1], arrLine26[2],
			arrLine26[3], arrLine26[4], arrLine26[5], arrLine26[6], arrLine26[7]},
		Line27: [9]string{strnumber[26], arrLine27[0], arrLine27[1], arrLine27[2],
			arrLine27[3], arrLine27[4], arrLine27[5], arrLine27[6], arrLine27[7]},
		Line28: [9]string{strnumber[27], arrLine28[0], arrLine28[1], arrLine28[2],
			arrLine28[3], arrLine28[4], arrLine28[5], arrLine28[6], arrLine28[7]},
		Line29: [9]string{strnumber[28], arrLine29[0], arrLine29[1], arrLine29[2],
			arrLine29[3], arrLine29[4], arrLine29[5], arrLine29[6], arrLine29[7]},
		Line30: [9]string{strnumber[29], arrLine30[0], arrLine30[1], arrLine30[2],
			arrLine30[3], arrLine30[4], arrLine30[5], arrLine30[6], arrLine30[7]},
		Line31: [9]string{strnumber[30], arrLine31[0], arrLine31[1], arrLine31[2],
			arrLine31[3], arrLine31[4], arrLine31[5], arrLine31[6], arrLine31[7]},
		Line32: [9]string{strnumber[31], arrLine32[0], arrLine32[1], arrLine32[2],
			arrLine32[3], arrLine32[4], arrLine32[5], arrLine32[6], arrLine32[7]},
		Line33: [9]string{strnumber[32], arrLine33[0], arrLine33[1], arrLine33[2],
			arrLine33[3], arrLine33[4], arrLine33[5], arrLine33[6], arrLine33[7]},
		Line34: [9]string{strnumber[33], arrLine34[0], arrLine34[1], arrLine34[2],
			arrLine34[3], arrLine34[4], arrLine34[5], arrLine34[6], arrLine34[7]},
		Line35: [9]string{strnumber[34], arrLine35[0], arrLine35[1], arrLine35[2],
			arrLine35[3], arrLine35[4], arrLine35[5], arrLine35[6], arrLine35[7]},
		Line36: [9]string{strnumber[35], arrLine36[0], arrLine36[1], arrLine36[2],
			arrLine36[3], arrLine36[4], arrLine36[5], arrLine36[6], arrLine36[7]},
		Line37: [9]string{strnumber[36], arrLine37[0], arrLine37[1], arrLine37[2],
			arrLine37[3], arrLine37[4], arrLine37[5], arrLine37[6], arrLine37[7]},
		Line38: [9]string{strnumber[37], arrLine38[0], arrLine38[1], arrLine38[2],
			arrLine38[3], arrLine38[4], arrLine38[5], arrLine38[6], arrLine38[7]},
		Line39: [9]string{strnumber[38], arrLine39[0], arrLine39[1], arrLine39[2],
			arrLine39[3], arrLine39[4], arrLine39[5], arrLine39[6], arrLine39[7]},
		Line40: [9]string{strnumber[39], arrLine40[0], arrLine40[1], arrLine40[2],
			arrLine40[3], arrLine40[4], arrLine40[5], arrLine40[6], arrLine40[7]},
		Line41: [9]string{strnumber[40], arrLine41[0], arrLine41[1], arrLine41[2],
			arrLine41[3], arrLine41[4], arrLine41[5], arrLine41[6], arrLine41[7]},
		Line42: [9]string{strnumber[41], arrLine42[0], arrLine42[1], arrLine42[2],
			arrLine42[3], arrLine42[4], arrLine42[5], arrLine42[6], arrLine42[7]},
		Line43: [9]string{strnumber[42], arrLine43[0], arrLine43[1], arrLine43[2],
			arrLine43[3], arrLine43[4], arrLine43[5], arrLine43[6], arrLine43[7]},
		Line44: [9]string{strnumber[43], arrLine44[0], arrLine44[1], arrLine44[2],
			arrLine44[3], arrLine44[4], arrLine44[5], arrLine44[6], arrLine44[7]},
		Line45: [9]string{strnumber[44], arrLine45[0], arrLine45[1], arrLine45[2],
			arrLine45[3], arrLine45[4], arrLine45[5], arrLine45[6], arrLine45[7]},
		Line46: [9]string{strnumber[45], arrLine46[0], arrLine46[1], arrLine46[2],
			arrLine46[3], arrLine46[4], arrLine46[5], arrLine46[6], arrLine46[7]},
		Line47: [9]string{strnumber[46], arrLine47[0], arrLine47[1], arrLine47[2],
			arrLine47[3], arrLine47[4], arrLine47[5], arrLine47[6], arrLine47[7]},
		Line48: [9]string{strnumber[47], arrLine48[0], arrLine48[1], arrLine48[2],
			arrLine48[3], arrLine48[4], arrLine48[5], arrLine48[6], arrLine48[7]},
		Line49: [9]string{strnumber[48], arrLine49[0], arrLine49[1], arrLine49[2],
			arrLine49[3], arrLine49[4], arrLine49[5], arrLine49[6], arrLine49[7]},
		Line50: [9]string{strnumber[49], arrLine50[0], arrLine50[1], arrLine50[2],
			arrLine50[3], arrLine50[4], arrLine50[5], arrLine50[6], arrLine50[7]},
		Line51: [9]string{strnumber[50], arrLine51[0], arrLine51[1], arrLine51[2],
			arrLine51[3], arrLine51[4], arrLine51[5], arrLine51[6], arrLine51[7]},
		Line52: [9]string{strnumber[51], arrLine52[0], arrLine52[1], arrLine52[2],
			arrLine52[3], arrLine52[4], arrLine52[5], arrLine52[6], arrLine52[7]},
		Line53: [9]string{strnumber[52], arrLine53[0], arrLine53[1], arrLine53[2],
			arrLine53[3], arrLine53[4], arrLine53[5], arrLine53[6], arrLine53[7]},
		Line54: [9]string{strnumber[53], arrLine54[0], arrLine54[1], arrLine54[2],
			arrLine54[3], arrLine54[4], arrLine54[5], arrLine54[6], arrLine54[7]},
		Line55: [9]string{strnumber[54], arrLine55[0], arrLine55[1], arrLine55[2],
			arrLine55[3], arrLine55[4], arrLine55[5], arrLine55[6], arrLine55[7]},
		Line56: [9]string{strnumber[55], arrLine56[0], arrLine56[1], arrLine56[2],
			arrLine56[3], arrLine56[4], arrLine56[5], arrLine56[6], arrLine56[7]},
		Line57: [9]string{strnumber[56], arrLine57[0], arrLine57[1], arrLine57[2],
			arrLine57[3], arrLine57[4], arrLine57[5], arrLine57[6], arrLine57[7]},
		Line58: [9]string{strnumber[57], arrLine58[0], arrLine58[1], arrLine58[2],
			arrLine58[3], arrLine58[4], arrLine58[5], arrLine58[6], arrLine58[7]},
		Line59: [9]string{strnumber[58], arrLine59[0], arrLine59[1], arrLine59[2],
			arrLine59[3], arrLine59[4], arrLine59[5], arrLine59[6], arrLine59[7]},
		Line60: [9]string{strnumber[59], arrLine60[0], arrLine60[1], arrLine60[2],
			arrLine60[3], arrLine60[4], arrLine60[5], arrLine60[6], arrLine60[7]},
	}
}

func bufferResult(i int) {
	switch i {
	case 0:
		arrLine1[0] = result.Tanggal
		arrLine1[1] = result.Jam
		arrLine1[2] = result.Heading
		arrLine1[3] = result.Latitude
		arrLine1[4] = result.Longitude
		arrLine1[5] = result.Speed
		arrLine1[6] = result.Temperature
		arrLine1[7] = result.CO2
	case 1:
		arrLine2[0] = result.Tanggal
		arrLine2[1] = result.Jam
		arrLine2[2] = result.Heading
		arrLine2[3] = result.Latitude
		arrLine2[4] = result.Longitude
		arrLine2[5] = result.Speed
		arrLine2[6] = result.Temperature
		arrLine2[7] = result.CO2
	case 2:
		arrLine3[0] = result.Tanggal
		arrLine3[1] = result.Jam
		arrLine3[2] = result.Heading
		arrLine3[3] = result.Latitude
		arrLine3[4] = result.Longitude
		arrLine3[5] = result.Speed
		arrLine3[6] = result.Temperature
		arrLine3[7] = result.CO2
	case 3:
		arrLine4[0] = result.Tanggal
		arrLine4[1] = result.Jam
		arrLine4[2] = result.Heading
		arrLine4[3] = result.Latitude
		arrLine4[4] = result.Longitude
		arrLine4[5] = result.Speed
		arrLine4[6] = result.Temperature
		arrLine4[7] = result.CO2
	case 4:
		arrLine5[0] = result.Tanggal
		arrLine5[1] = result.Jam
		arrLine5[2] = result.Heading
		arrLine5[3] = result.Latitude
		arrLine5[4] = result.Longitude
		arrLine5[5] = result.Speed
		arrLine5[6] = result.Temperature
		arrLine5[7] = result.CO2
	case 5:
		arrLine6[0] = result.Tanggal
		arrLine6[1] = result.Jam
		arrLine6[2] = result.Heading
		arrLine6[3] = result.Latitude
		arrLine6[4] = result.Longitude
		arrLine6[5] = result.Speed
		arrLine6[6] = result.Temperature
		arrLine6[7] = result.CO2
	case 6:
		arrLine7[0] = result.Tanggal
		arrLine7[1] = result.Jam
		arrLine7[2] = result.Heading
		arrLine7[3] = result.Latitude
		arrLine7[4] = result.Longitude
		arrLine7[5] = result.Speed
		arrLine7[6] = result.Temperature
		arrLine7[7] = result.CO2
	case 7:
		arrLine8[0] = result.Tanggal
		arrLine8[1] = result.Jam
		arrLine8[2] = result.Heading
		arrLine8[3] = result.Latitude
		arrLine8[4] = result.Longitude
		arrLine8[5] = result.Speed
		arrLine8[6] = result.Temperature
		arrLine8[7] = result.CO2
	case 8:
		arrLine9[0] = result.Tanggal
		arrLine9[1] = result.Jam
		arrLine9[2] = result.Heading
		arrLine9[3] = result.Latitude
		arrLine9[4] = result.Longitude
		arrLine9[5] = result.Speed
		arrLine9[6] = result.Temperature
		arrLine9[7] = result.CO2
	case 9:
		arrLine10[0] = result.Tanggal
		arrLine10[1] = result.Jam
		arrLine10[2] = result.Heading
		arrLine10[3] = result.Latitude
		arrLine10[4] = result.Longitude
		arrLine10[5] = result.Speed
		arrLine10[6] = result.Temperature
		arrLine10[7] = result.CO2
	case 10:
		arrLine11[0] = result.Tanggal
		arrLine11[1] = result.Jam
		arrLine11[2] = result.Heading
		arrLine11[3] = result.Latitude
		arrLine11[4] = result.Longitude
		arrLine11[5] = result.Speed
		arrLine11[6] = result.Temperature
		arrLine11[7] = result.CO2
	case 11:
		arrLine12[0] = result.Tanggal
		arrLine12[1] = result.Jam
		arrLine12[2] = result.Heading
		arrLine12[3] = result.Latitude
		arrLine12[4] = result.Longitude
		arrLine12[5] = result.Speed
		arrLine12[6] = result.Temperature
		arrLine12[7] = result.CO2
	case 12:
		arrLine13[0] = result.Tanggal
		arrLine13[1] = result.Jam
		arrLine13[2] = result.Heading
		arrLine13[3] = result.Latitude
		arrLine13[4] = result.Longitude
		arrLine13[5] = result.Speed
		arrLine13[6] = result.Temperature
		arrLine13[7] = result.CO2
	case 13:
		arrLine14[0] = result.Tanggal
		arrLine14[1] = result.Jam
		arrLine14[2] = result.Heading
		arrLine14[3] = result.Latitude
		arrLine14[4] = result.Longitude
		arrLine14[5] = result.Speed
		arrLine14[6] = result.Temperature
		arrLine14[7] = result.CO2
	case 14:
		arrLine15[0] = result.Tanggal
		arrLine15[1] = result.Jam
		arrLine15[2] = result.Heading
		arrLine15[3] = result.Latitude
		arrLine15[4] = result.Longitude
		arrLine15[5] = result.Speed
		arrLine15[6] = result.Temperature
		arrLine15[7] = result.CO2
	case 15:
		arrLine16[0] = result.Tanggal
		arrLine16[1] = result.Jam
		arrLine16[2] = result.Heading
		arrLine16[3] = result.Latitude
		arrLine16[4] = result.Longitude
		arrLine16[5] = result.Speed
		arrLine16[6] = result.Temperature
		arrLine16[7] = result.CO2
	case 16:
		arrLine17[0] = result.Tanggal
		arrLine17[1] = result.Jam
		arrLine17[2] = result.Heading
		arrLine17[3] = result.Latitude
		arrLine17[4] = result.Longitude
		arrLine17[5] = result.Speed
		arrLine17[6] = result.Temperature
		arrLine17[7] = result.CO2
	case 17:
		arrLine18[0] = result.Tanggal
		arrLine18[1] = result.Jam
		arrLine18[2] = result.Heading
		arrLine18[3] = result.Latitude
		arrLine18[4] = result.Longitude
		arrLine18[5] = result.Speed
		arrLine18[6] = result.Temperature
		arrLine18[7] = result.CO2
	case 18:
		arrLine19[0] = result.Tanggal
		arrLine19[1] = result.Jam
		arrLine19[2] = result.Heading
		arrLine19[3] = result.Latitude
		arrLine19[4] = result.Longitude
		arrLine19[5] = result.Speed
		arrLine19[6] = result.Temperature
		arrLine19[7] = result.CO2
	case 19:
		arrLine20[0] = result.Tanggal
		arrLine20[1] = result.Jam
		arrLine20[2] = result.Heading
		arrLine20[3] = result.Latitude
		arrLine20[4] = result.Longitude
		arrLine20[5] = result.Speed
		arrLine20[6] = result.Temperature
		arrLine20[7] = result.CO2
	case 20:
		arrLine21[0] = result.Tanggal
		arrLine21[1] = result.Jam
		arrLine21[2] = result.Heading
		arrLine21[3] = result.Latitude
		arrLine21[4] = result.Longitude
		arrLine21[5] = result.Speed
		arrLine21[6] = result.Temperature
		arrLine21[7] = result.CO2
	case 21:
		arrLine22[0] = result.Tanggal
		arrLine22[1] = result.Jam
		arrLine22[2] = result.Heading
		arrLine22[3] = result.Latitude
		arrLine22[4] = result.Longitude
		arrLine22[5] = result.Speed
		arrLine22[6] = result.Temperature
		arrLine22[7] = result.CO2
	case 22:
		arrLine23[0] = result.Tanggal
		arrLine23[1] = result.Jam
		arrLine23[2] = result.Heading
		arrLine23[3] = result.Latitude
		arrLine23[4] = result.Longitude
		arrLine23[5] = result.Speed
		arrLine23[6] = result.Temperature
		arrLine23[7] = result.CO2
	case 23:
		arrLine24[0] = result.Tanggal
		arrLine24[1] = result.Jam
		arrLine24[2] = result.Heading
		arrLine24[3] = result.Latitude
		arrLine24[4] = result.Longitude
		arrLine24[5] = result.Speed
		arrLine24[6] = result.Temperature
		arrLine24[7] = result.CO2
	case 24:
		arrLine25[0] = result.Tanggal
		arrLine25[1] = result.Jam
		arrLine25[2] = result.Heading
		arrLine25[3] = result.Latitude
		arrLine25[4] = result.Longitude
		arrLine25[5] = result.Speed
		arrLine25[6] = result.Temperature
		arrLine25[7] = result.CO2
	case 25:
		arrLine26[0] = result.Tanggal
		arrLine26[1] = result.Jam
		arrLine26[2] = result.Heading
		arrLine26[3] = result.Latitude
		arrLine26[4] = result.Longitude
		arrLine26[5] = result.Speed
		arrLine26[6] = result.Temperature
		arrLine26[7] = result.CO2
	case 26:
		arrLine27[0] = result.Tanggal
		arrLine27[1] = result.Jam
		arrLine27[2] = result.Heading
		arrLine27[3] = result.Latitude
		arrLine27[4] = result.Longitude
		arrLine27[5] = result.Speed
		arrLine27[6] = result.Temperature
		arrLine27[7] = result.CO2
	case 27:
		arrLine28[0] = result.Tanggal
		arrLine28[1] = result.Jam
		arrLine28[2] = result.Heading
		arrLine28[3] = result.Latitude
		arrLine28[4] = result.Longitude
		arrLine28[5] = result.Speed
		arrLine28[6] = result.Temperature
		arrLine28[7] = result.CO2
	case 28:
		arrLine29[0] = result.Tanggal
		arrLine29[1] = result.Jam
		arrLine29[2] = result.Heading
		arrLine29[3] = result.Latitude
		arrLine29[4] = result.Longitude
		arrLine29[5] = result.Speed
		arrLine29[6] = result.Temperature
		arrLine29[7] = result.CO2
	case 29:
		arrLine30[0] = result.Tanggal
		arrLine30[1] = result.Jam
		arrLine30[2] = result.Heading
		arrLine30[3] = result.Latitude
		arrLine30[4] = result.Longitude
		arrLine30[5] = result.Speed
		arrLine30[6] = result.Temperature
		arrLine30[7] = result.CO2
	case 30:
		arrLine31[0] = result.Tanggal
		arrLine31[1] = result.Jam
		arrLine31[2] = result.Heading
		arrLine31[3] = result.Latitude
		arrLine31[4] = result.Longitude
		arrLine31[5] = result.Speed
		arrLine31[6] = result.Temperature
		arrLine31[7] = result.CO2
	case 31:
		arrLine32[0] = result.Tanggal
		arrLine32[1] = result.Jam
		arrLine32[2] = result.Heading
		arrLine32[3] = result.Latitude
		arrLine32[4] = result.Longitude
		arrLine32[5] = result.Speed
		arrLine32[6] = result.Temperature
		arrLine32[7] = result.CO2
	case 32:
		arrLine33[0] = result.Tanggal
		arrLine33[1] = result.Jam
		arrLine33[2] = result.Heading
		arrLine33[3] = result.Latitude
		arrLine33[4] = result.Longitude
		arrLine33[5] = result.Speed
		arrLine33[6] = result.Temperature
		arrLine33[7] = result.CO2
	case 33:
		arrLine34[0] = result.Tanggal
		arrLine34[1] = result.Jam
		arrLine34[2] = result.Heading
		arrLine34[3] = result.Latitude
		arrLine34[4] = result.Longitude
		arrLine34[5] = result.Speed
		arrLine34[6] = result.Temperature
		arrLine34[7] = result.CO2
	case 34:
		arrLine35[0] = result.Tanggal
		arrLine35[1] = result.Jam
		arrLine35[2] = result.Heading
		arrLine35[3] = result.Latitude
		arrLine35[4] = result.Longitude
		arrLine35[5] = result.Speed
		arrLine35[6] = result.Temperature
		arrLine35[7] = result.CO2
	case 35:
		arrLine36[0] = result.Tanggal
		arrLine36[1] = result.Jam
		arrLine36[2] = result.Heading
		arrLine36[3] = result.Latitude
		arrLine36[4] = result.Longitude
		arrLine36[5] = result.Speed
		arrLine36[6] = result.Temperature
		arrLine36[7] = result.CO2
	case 36:
		arrLine37[0] = result.Tanggal
		arrLine37[1] = result.Jam
		arrLine37[2] = result.Heading
		arrLine37[3] = result.Latitude
		arrLine37[4] = result.Longitude
		arrLine37[5] = result.Speed
		arrLine37[6] = result.Temperature
		arrLine37[7] = result.CO2
	case 37:
		arrLine38[0] = result.Tanggal
		arrLine38[1] = result.Jam
		arrLine38[2] = result.Heading
		arrLine38[3] = result.Latitude
		arrLine38[4] = result.Longitude
		arrLine38[5] = result.Speed
		arrLine38[6] = result.Temperature
		arrLine38[7] = result.CO2
	case 38:
		arrLine39[0] = result.Tanggal
		arrLine39[1] = result.Jam
		arrLine39[2] = result.Heading
		arrLine39[3] = result.Latitude
		arrLine39[4] = result.Longitude
		arrLine39[5] = result.Speed
		arrLine39[6] = result.Temperature
		arrLine39[7] = result.CO2
	case 39:
		arrLine40[0] = result.Tanggal
		arrLine40[1] = result.Jam
		arrLine40[2] = result.Heading
		arrLine40[3] = result.Latitude
		arrLine40[4] = result.Longitude
		arrLine40[5] = result.Speed
		arrLine40[6] = result.Temperature
		arrLine40[7] = result.CO2

	case 40:
		arrLine41[0] = result.Tanggal
		arrLine41[1] = result.Jam
		arrLine41[2] = result.Heading
		arrLine41[3] = result.Latitude
		arrLine41[4] = result.Longitude
		arrLine41[5] = result.Speed
		arrLine41[6] = result.Temperature
		arrLine41[7] = result.CO2
	case 41:
		arrLine42[0] = result.Tanggal
		arrLine42[1] = result.Jam
		arrLine42[2] = result.Heading
		arrLine42[3] = result.Latitude
		arrLine42[4] = result.Longitude
		arrLine42[5] = result.Speed
		arrLine42[6] = result.Temperature
		arrLine42[7] = result.CO2
	case 42:
		arrLine43[0] = result.Tanggal
		arrLine43[1] = result.Jam
		arrLine43[2] = result.Heading
		arrLine43[3] = result.Latitude
		arrLine43[4] = result.Longitude
		arrLine43[5] = result.Speed
		arrLine43[6] = result.Temperature
		arrLine43[7] = result.CO2
	case 43:
		arrLine44[0] = result.Tanggal
		arrLine44[1] = result.Jam
		arrLine44[2] = result.Heading
		arrLine44[3] = result.Latitude
		arrLine44[4] = result.Longitude
		arrLine44[5] = result.Speed
		arrLine44[6] = result.Temperature
		arrLine44[7] = result.CO2
	case 44:
		arrLine45[0] = result.Tanggal
		arrLine45[1] = result.Jam
		arrLine45[2] = result.Heading
		arrLine45[3] = result.Latitude
		arrLine45[4] = result.Longitude
		arrLine45[5] = result.Speed
		arrLine45[6] = result.Temperature
		arrLine45[7] = result.CO2
	case 45:
		arrLine46[0] = result.Tanggal
		arrLine46[1] = result.Jam
		arrLine46[2] = result.Heading
		arrLine46[3] = result.Latitude
		arrLine46[4] = result.Longitude
		arrLine46[5] = result.Speed
		arrLine46[6] = result.Temperature
		arrLine46[7] = result.CO2
	case 46:
		arrLine47[0] = result.Tanggal
		arrLine47[1] = result.Jam
		arrLine47[2] = result.Heading
		arrLine47[3] = result.Latitude
		arrLine47[4] = result.Longitude
		arrLine47[5] = result.Speed
		arrLine47[6] = result.Temperature
		arrLine47[7] = result.CO2
	case 47:
		arrLine48[0] = result.Tanggal
		arrLine48[1] = result.Jam
		arrLine48[2] = result.Heading
		arrLine48[3] = result.Latitude
		arrLine48[4] = result.Longitude
		arrLine48[5] = result.Speed
		arrLine48[6] = result.Temperature
		arrLine48[7] = result.CO2
	case 48:
		arrLine49[0] = result.Tanggal
		arrLine49[1] = result.Jam
		arrLine49[2] = result.Heading
		arrLine49[3] = result.Latitude
		arrLine49[4] = result.Longitude
		arrLine49[5] = result.Speed
		arrLine49[6] = result.Temperature
		arrLine49[7] = result.CO2
	case 49:
		arrLine50[0] = result.Tanggal
		arrLine50[1] = result.Jam
		arrLine50[2] = result.Heading
		arrLine50[3] = result.Latitude
		arrLine50[4] = result.Longitude
		arrLine50[5] = result.Speed
		arrLine50[6] = result.Temperature
		arrLine50[7] = result.CO2

	case 50:
		arrLine51[0] = result.Tanggal
		arrLine51[1] = result.Jam
		arrLine51[2] = result.Heading
		arrLine51[3] = result.Latitude
		arrLine51[4] = result.Longitude
		arrLine51[5] = result.Speed
		arrLine51[6] = result.Temperature
		arrLine51[7] = result.CO2
	case 51:
		arrLine52[0] = result.Tanggal
		arrLine52[1] = result.Jam
		arrLine52[2] = result.Heading
		arrLine52[3] = result.Latitude
		arrLine52[4] = result.Longitude
		arrLine52[5] = result.Speed
		arrLine52[6] = result.Temperature
		arrLine52[7] = result.CO2
	case 52:
		arrLine53[0] = result.Tanggal
		arrLine53[1] = result.Jam
		arrLine53[2] = result.Heading
		arrLine53[3] = result.Latitude
		arrLine53[4] = result.Longitude
		arrLine53[5] = result.Speed
		arrLine53[6] = result.Temperature
		arrLine53[7] = result.CO2
	case 53:
		arrLine54[0] = result.Tanggal
		arrLine54[1] = result.Jam
		arrLine54[2] = result.Heading
		arrLine54[3] = result.Latitude
		arrLine54[4] = result.Longitude
		arrLine54[5] = result.Speed
		arrLine54[6] = result.Temperature
		arrLine54[7] = result.CO2
	case 54:
		arrLine55[0] = result.Tanggal
		arrLine55[1] = result.Jam
		arrLine55[2] = result.Heading
		arrLine55[3] = result.Latitude
		arrLine55[4] = result.Longitude
		arrLine55[5] = result.Speed
		arrLine55[6] = result.Temperature
		arrLine55[7] = result.CO2
	case 55:
		arrLine56[0] = result.Tanggal
		arrLine56[1] = result.Jam
		arrLine56[2] = result.Heading
		arrLine56[3] = result.Latitude
		arrLine56[4] = result.Longitude
		arrLine56[5] = result.Speed
		arrLine56[6] = result.Temperature
		arrLine56[7] = result.CO2
	case 56:
		arrLine57[0] = result.Tanggal
		arrLine57[1] = result.Jam
		arrLine57[2] = result.Heading
		arrLine57[3] = result.Latitude
		arrLine57[4] = result.Longitude
		arrLine57[5] = result.Speed
		arrLine57[6] = result.Temperature
		arrLine57[7] = result.CO2
	case 57:
		arrLine58[0] = result.Tanggal
		arrLine58[1] = result.Jam
		arrLine58[2] = result.Heading
		arrLine58[3] = result.Latitude
		arrLine58[4] = result.Longitude
		arrLine58[5] = result.Speed
		arrLine58[6] = result.Temperature
		arrLine58[7] = result.CO2
	case 58:
		arrLine59[0] = result.Tanggal
		arrLine59[1] = result.Jam
		arrLine59[2] = result.Heading
		arrLine59[3] = result.Latitude
		arrLine59[4] = result.Longitude
		arrLine59[5] = result.Speed
		arrLine59[6] = result.Temperature
		arrLine59[7] = result.CO2
	case 59:
		arrLine60[0] = result.Tanggal
		arrLine60[1] = result.Jam
		arrLine60[2] = result.Heading
		arrLine60[3] = result.Latitude
		arrLine60[4] = result.Longitude
		arrLine60[5] = result.Speed
		arrLine60[6] = result.Temperature
		arrLine60[7] = result.CO2

	default:
		break
	}
}

func routemap17(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("map17").ParseFiles("view/map.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		strmin17 := strconv.Itoa(min)
		strmax17 := strconv.Itoa(max)
		var datamap = map[string]string{"jamin": jamin17, "datemin": datemin17, "latmin": latmin17, "longmin": longmin17, "comin": strmin17, "tempmin": tempmin17, "speedmin": speedmin17,
			"jamax": jamax17, "datemax": datemax17, "latmax": latmax17, "longmax": longmax17, "comax": strmax17, "tempmax": tempmax17, "speedmax": speedmax17}

		if err := tmpl.Execute(w, datamap); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routemap29(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("map29").ParseFiles("view/map.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		strmin := strconv.Itoa(min)
		strmax := strconv.Itoa(max)
		var datamap = map[string]string{"jamin": jamin, "datemin": datemin, "latmin": latmin, "longmin": longmin, "comin": strmin, "tempmin": tempmin, "speedmin": speedmin,
			"jamax": jamax, "datemax": datemax, "latmax": latmax, "longmax": longmax, "comax": strmax, "tempmax": tempmax, "speedmax": speedmax}

		if err := tmpl.Execute(w, datamap); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routehome(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("home").ParseFiles("view/home.html"))

		var err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
