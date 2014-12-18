package main

import "github.com/docopt/docopt-go"
import . "github.com/tj/go-debug"
import "strings"
import "strconv"
import "net/url"
import "fmt"
import "os"

var method, uri, headers, body string
var limitRequestTimes int
var totalRequestTimes int
var co, recoverTimes int

var debug = Debug("httpbench")

func main() {
	usage := `
	Usage:
		httpbench [-u=<url>] [-m=<method>] [-c=<concurrent>] [-h=<headers>] [-b=<body>] [-l=<limit>]
		httpbench --help
		httpbench --version

	Options:
		-u=<url>        Required, url to bench
		-m=<method>     Add method, such as: GET
		-c=<concurrent> Set number of requests to run concurrently
		-h=<headers>    Add headers, such as: "Content-Type:text/xml; Content-Length:100"
		-b=<body>       Add body, such as: "name=haoxin&age=24"
		-l=<limit>      Set limit for request times
		--help          Show this screen
		--version       Show version
	`

	args, _ := docopt.Parse(usage, os.Args[1:], true, "v0.1.0", false)
	debug("args: %v", args)
	parse(args)

	quit := make(chan bool)

	for i := 0; i < co; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					recoverTimes++
					fmt.Printf("panic %d in main, message: %v \n", recoverTimes, r)
					if recoverTimes >= co {
						os.Exit(1)
					}
				}
			}()

			for {
				request(method, uri, headers, body)

				if totalRequestTimes >= limitRequestTimes {
					fmt.Println("reach limit times")
					os.Exit(1)
				}
			}
		}()
	}

	if <-quit {
	}
}

func parse(args map[string]interface{}) {
	for k, v := range args {
		switch k {
		case "-h":
			if v != nil {
				headers = v.(string)
			}
		case "-b":
			if v != nil {
				body = v.(string)
			}
		case "-m":
			var m string
			if v != nil {
				m = strings.ToUpper(v.(string))
			}
			if m == "" {
				method = "GET"
			} else {
				method = m
			}
		case "-u":
			uri = getUrl(v)
		case "-c":
			co = getInt(v)
		case "-l":
			limitRequestTimes = getInt(v)
		}
	}

	if uri == "" {
		fmt.Println(`request url is required, use -u "your url"`)
		os.Exit(1)
	}

	if method == "" {
		method = "GET"
	}

	if co <= 0 {
		co = 5
	}

	if limitRequestTimes <= 0 {
		limitRequestTimes = 10000
	}

	fmt.Printf("concurrency: %d, method: %s, url: %s \n", co, method, uri)
	if headers != "" {
		fmt.Printf("headers: %s \n", headers)
	}
	if body != "" {
		fmt.Printf("body: %s \n", body)
	}
}

func getUrl(i interface{}) string {
	if i == nil {
		return ""
	}

	s := i.(string)
	var uri string
	if !strings.HasPrefix(s, "http") {
		uri = "http://" + s
	} else {
		uri = s
	}

	u, e := url.ParseRequestURI(uri)
	debug("URL: %v, err: %v", u, e)
	if e != nil {
		fmt.Println("invalid url")
		os.Exit(1)
	}
	return u.String()
}

func getInt(i interface{}) int {
	if i == nil {
		return 0
	}

	s := i.(string)
	num, e := strconv.Atoi(s)

	if e != nil {
		return 0
	}

	return num
}
