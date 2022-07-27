package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://google.com",
		"https://facebook.com",
		"https://stackoverflow.com",
		"https://golang.org",
		"https://amazon.com",
	}
	fmt.Println("WaitGroup **************")
	waitGroup(urls)
	fmt.Println("")
	fmt.Println("Channels ***************")
	channels(urls)
	fmt.Println("")
	fmt.Println("Channels infinite ******")
	channelsInfinite(urls)
}

// channels sol
func channels(urls []string) {
	c := make(chan string)
	for _, url := range urls {
		go verifyUrlIsUpWithChannel(url, c)
	}
	for i := 0; i < len(urls); i++ {
		fmt.Println(<-c)
	}
}

// channels infinite
func channelsInfinite(urls []string) {
	c := make(chan string)
	for _, url := range urls {
		go verifyUrlIsUpWithChannel(url, c)
	}
	// option A
	// for {
	// 	go verifyUrlIsUpWithChannel(<-c, c)
	// }

	// option B: more readable
	for link := range c {
		// we need to forward l because we need it to be passed as value or as a copy,
		// otherwise the link value could be mutated by the main goroutine before the verifyUrlIsUpWithChannel is executed
		// TIP: never share between go routines data. always forward args
		go func(l string) {
			time.Sleep(5 * time.Second)
			verifyUrlIsUpWithChannel(l, c)
		}(link)
	}
}

func verifyUrlIsUpWithChannel(url string, c chan string) {
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(url, "is down")
	} else {
		fmt.Println(url, "is up")
	}
	c <- url
}

// wait group sol
func waitGroup(urls []string) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go verifyUrlIsUp(url, &wg)
	}
	wg.Wait()
}

func verifyUrlIsUp(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(url, " is down")
	} else {
		fmt.Println(url, " is up")
	}
}
