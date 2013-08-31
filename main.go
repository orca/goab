// goab project main.go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var l sync.Mutex
var finished = 0
var sig = make(chan int)
var error_request = 0
var success_request = 0

func main() {
	c := flag.Int("c", 10, "Concurrency level")
	n := flag.Int("n", 1000, "otal requests")
	url := flag.String("url", "http://localhost", "The url to test benchmark")

	flag.Parse()

	pool := make(chan string, *c)

	time_begin := time.Now()
	go initPool(pool, *n, *url)
	for i := 0; i < *c; i++ {
		go requestUrl(pool, *n)
	}
	<-sig
	t := time.Now().Sub(time_begin)

	fmt.Println("Total time", t.Seconds())
	fmt.Println("success request:", success_request)
	fmt.Println("error request:", error_request)
	fmt.Println("Time per request", t.Seconds()/float64(*n))
	fmt.Println("benchmark:", float64(*n)/t.Seconds())
	fmt.Println("n", *n)
	fmt.Println("c", *c)
	fmt.Println("url", *url)
}

func initPool(pool chan string, n int, url string) {
	for i := 0; i < n; i++ {
		pool <- url
	}
}

func requestUrl(pool chan string, n int) {
	for {
		url := <-pool
		_, err := http.Get(url)
		l.Lock()
		if err != nil {
			error_request++
		} else {
			success_request++
		}
		finished++
		l.Unlock()
		if finished >= n {
			sig <- 1
		}
	}
}
