package main

import "io/ioutil"
import "net/http"
import "strings"
import "time"
import "fmt"

var client *http.Client
var concurrent = 0

func init() {
	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  true,
	}

	client = &http.Client{Transport: transport}
}

func request(method, httpUrl, headers, body string) {
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

	fmt.Printf("* response * length: %d, duration: %d ms, concurrent: %d \n", bodyLength, (end-start)/1000000, concurrent)

	panicError(err)
}

func panicError(e error) {
	if e != nil {
		panic(e)
	}
}
