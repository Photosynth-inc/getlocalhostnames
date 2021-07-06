package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	DIG_SLEEP = 1000
	QUEUE_WAIT = 200
)

func dig(ip string, verbose bool) {
	out, err := exec.Command("dig", "+short", "+time=1", "+tries=1", "-x", ip, "@224.0.0.251", "-p", "5353").Output()
	if err != nil && verbose {
		fmt.Printf("%s\tFailed: %s\n", ip, err)
		return
	} 

	fmt.Printf("%s\t%s", ip, string(out))
	time.Sleep(DIG_SLEEP * time.Millisecond)
}

func main() {
	limit := make(chan struct{}, 50)

	verbose := false
	for _, v := range os.Args {
		switch v {
			case "-v" :
				verbose = true
		}
	}

	var wg sync.WaitGroup
	for i := 1; i < 255; i++ {
		ip := fmt.Sprintf("192.168.101.%d", i)

		wg.Add(1)
		go func() {
			limit <- struct{}{}
			defer wg.Done()
			dig(ip, verbose)
			<-limit
		}()
		time.Sleep(QUEUE_WAIT * time.Millisecond)
	}
	wg.Wait()
}
