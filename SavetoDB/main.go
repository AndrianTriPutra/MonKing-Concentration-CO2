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
	fmt.Println("Fimrware Version:29/10/2019")
	setup()
	loop()
}

func loop() {
	var wg sync.WaitGroup
	stopchan := make(chan bool, 1)

	go func() {
		wg.Add(1)
		defer wg.Done()

		tickersubscriber := time.NewTicker(10 * time.Second)
		defer tickersubscriber.Stop()

		subscriber(*tickersubscriber, stopchan)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		generate(stopchan)
	}()

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, syscall.SIGINT, syscall.SIGTERM)
	<-cc

	close(stopchan)
	wg.Wait()
}
