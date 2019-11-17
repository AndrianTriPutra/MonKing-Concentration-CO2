package main

import (
	"fmt"
	"net/http"
	"time"
)

func nethttp() {
	http.HandleFunc("/monitoring/co2", handlerco2)

	var address = "localhost:9000"
	fmt.Printf("server started at %s\n", address)

	server := new(http.Server)
	server.Addr = address
	server.ReadTimeout = time.Second * 5
	server.WriteTimeout = time.Second * 5
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func handlerco2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(printout))
}
