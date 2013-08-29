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

func main() {
	c := flag.Int("c", 10, "The number of concurrent")
	n := flag.Int("n", 1000, "The number of total requests")
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

	fmt.Println("Total time Nano", t.Nanoseconds())
	fmt.Println("Total time second", t.Seconds())
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
		if err != nil {

		} else {

		}

		l.Lock()
		finished = finished + 1
		l.Unlock()
		if finished >= n {
			sig <- 1
		}
	}
}
