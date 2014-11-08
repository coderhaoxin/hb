package main

import "io/ioutil"
import "net/http"
import "strings"
import "time"
import "fmt"

var client *http.Client
var total = 0

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

	if err != nil {
	}

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
	total++
	start := time.Now().UnixNano()
	res, err := client.Do(req)

	if err != nil {
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	end := time.Now().UnixNano()
	total--

	bodyLength := len(data)

	fmt.Printf("* response * length: %d, duration: %d ms, total: %d \n", bodyLength, (end-start)/1000000, total)

	if err != nil {
	}
}
