package main

import "strings"
import "flag"
import "fmt"
import "os"

var headers = flag.String("h", "", "http request headers, such as: 'Content-Type:text/xml; Content-Length:100'")
var body = flag.String("b", "", "http request body, such as: 'name=haoxin&age=24'")
var method = flag.String("m", "", "http request method, such as: GET")
var url = flag.String("u", "", "http request url")
var co = flag.Int("c", 10, "number of requests to run concurrently")

var m, u string

func main() {
	flag.Parse()
	assert()

	quit := make(chan bool)

	for i := 0; i < *co; i++ {
		go func() {
			for {
				request(m, u, *headers, *body)
			}
		}()
	}

	if <-quit {
	}
}

func assert() {
	if *url == "" {
		fmt.Println("request url is required, use -u 'your url'")
		os.Exit(1)
	}

	if strings.ToUpper(*method) != "" {
		m = strings.ToUpper(*method)
	} else {
		m = "GET"
	}

	if strings.HasSuffix(*url, "http") {
		u = *url
	} else {
		u = "http://" + *url
	}

	fmt.Printf("concurrency: %d, method: %s, url: %s \n", *co, m, u)
	if *headers != "" {
		fmt.Printf("headers: %s \n", *headers)
	}
	if *body != "" {
		fmt.Printf("body: %s \n", *body)
	}
}
