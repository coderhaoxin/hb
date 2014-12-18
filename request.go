package main

import "io/ioutil"
import "net/http"
import "strings"
import "time"
import "fmt"

var client *http.Client
var concurrent = 0
var panicTimes int
var panics [10]int64

func init() {
	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  true,
	}

	client = &http.Client{Transport: transport}
}

func request(method, httpUrl, headers, body string) {
	defer catchPanic()

	reader := strings.NewReader(body)
	req, err := http.NewRequest(method, httpUrl, reader)

	panicError(err)

	// set headers
	for _, header := range strings.Split(headers, ";") {
		s := strings.Split(header, ":")

		if len(s) < 2 {
			continue
		}

		k := strings.TrimSpace(s[0])
		v := strings.TrimSpace(s[1])
		req.Header.Set(k, v)
	}

	// do http request
	concurrent++
	start := time.Now().UnixNano()
	res, err := client.Do(req)

	panicError(err)

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	end := time.Now().UnixNano()
	concurrent--

	bodyLength := len(data)

	totalRequestTimes++
	fmt.Printf("* response * length: %d, duration: %d ms, total %d, concurrent: %d \n", bodyLength, (end-start)/1000000, totalRequestTimes, concurrent)

	panicError(err)
}

func panicError(e error) {
	if e != nil {
		panic(e)
	}
}

func catchPanic() {
	end := len(panics) - 1
	panicTimes++

	for i := 0; i < end; i++ {
		panics[i] = panics[i+1]
	}

	panics[end] = time.Now().Unix()

	if panicTimes < 3*end || panics[end]-panics[0] > 3 {
		if r := recover(); r != nil {
			fmt.Printf("panic %d in request, message: %v \n", panicTimes, r)
		}
	}
}
