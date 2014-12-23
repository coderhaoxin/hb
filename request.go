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
var durations []int64

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
	duration := end - start
	durations = append(durations, duration)
	fmt.Printf("* response * length: %d, duration: %d ms, total %d, concurrent: %d \n", bodyLength, duration/1000000, totalRequestTimes, concurrent)

	panicError(err)
}

func panicError(e error) {
	if e != nil {
		panic(e)
	}
}

func catchPanic() {
	if r := recover(); r != nil {

		end := len(panics) - 1
		panicTimes++

		for i := 0; i < end; i++ {
			panics[i] = panics[i+1]
		}

		panics[end] = time.Now().Unix()

		if panicTimes > 3*end || panics[end]-panics[0] < 3 {
			fmt.Printf("panic %d in request, message: %v \n", panicTimes, r)
		}
	}
}

func analysis() (average int64, levels map[string]int64) {
	var total int64
	levels = make(map[string]int64)

	count := len(durations)

	for i := 0; i < count; i++ {
		total += durations[i] / 1000000
	}

	average = total / int64(count)
	level0 := int64(0)
	level1 := average / 5
	level2 := level1 * 2
	level3 := level1 * 3
	level4 := level1 * 4

	var counts [5]int64

	for i := 0; i < count; i++ {
		ms := durations[i] / 1000000
		switch {
		case ms > level4:
			counts[4]++
		case ms > level3:
			counts[3]++
		case ms > level2:
			counts[2]++
		case ms > level1:
			counts[1]++
		case ms > level0:
			counts[0]++
		}
	}

	levels[toString(level0)+"~"+toString(level1)] = level0
	levels[toString(level1)+"~"+toString(level2)] = level1
	levels[toString(level2)+"~"+toString(level3)] = level2
	levels[toString(level3)+"~"+toString(level4)] = level3
	levels["above "+toString(level4)] = level4

	return
}
