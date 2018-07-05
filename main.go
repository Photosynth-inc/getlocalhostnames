package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func doSomething(u int) {

	str := fmt.Sprintf("192.168.0.%d", u)
	//fmt.Println(str)
	out, err := exec.Command("dig", "+short", "+time=1", "+tries=1", "-x", str, "@224.0.0.251", "-p", "5353").Output()
	if err != nil {
		//log.Printf("Exec fail: %v", err)
	} else {
		fmt.Print(string(out))
	}
	// dig -x ip @224.0.0.251 -p 5353
	//time.Sleep(2 * time.Second)
}

func main() {

	limit := make(chan struct{}, 50)

	var wg sync.WaitGroup
	for i := 1; i < 255; i++ {
		ii := i
		wg.Add(1)
		go func() {
			limit <- struct{}{}
			defer wg.Done()
			doSomething(ii)
			<-limit
		}()
		//time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}
