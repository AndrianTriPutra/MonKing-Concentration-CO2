package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Monitoring CO2 29/10/19 \n")

	fmt.Println("warming up device")
	fmt.Println("       please wait . . .")
	time.Sleep(30 * time.Second)

	setup()

	var wg sync.WaitGroup
	stopchan := make(chan bool, 1)

	go func() {
		mqtt()
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		gps(*ticker, stopchan)

	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		getco2(*ticker, stopchan)

	}()

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, syscall.SIGINT, syscall.SIGTERM)
	<-cc

	close(stopchan)
	wg.Wait()

}
