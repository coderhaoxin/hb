package main

import "github.com/olekukonko/tablewriter"
import "github.com/docopt/docopt-go"
import . "github.com/tj/go-debug"
import "os/signal"
import "strings"
import "strconv"
import "net/url"
import "fmt"
import "os"

var method, uri, headers, body string
var limitRequestTimes int
var totalRequestTimes int
var co, recoverTimes int

var debug = Debug("hb")

const version = "v0.2.0"

func main() {
	usage := `
	Usage:
		hb [-u=<url>] [-m=<method>] [-c=<concurrent>] [--headers=<headers>] [--body=<body>] [--limit=<limit>]
		hb --help
		hb --version

	Options:
		-u=<url>            Required, url to bench
		-m=<method>         Add method, such as: GET
		-c=<concurrent>     Set number of requests to run concurrently
		--headers=<headers> Add headers, such as: "Content-Type:text/xml; Content-Length:100"
		--body=<body>       Add body, such as: "name=haoxin&age=24"
		--limit=<limit>     Set limit for request times
		--help              Show this screen
		--version           Show version
	`

	args, _ := docopt.Parse(usage, os.Args[1:], true, version, false)
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
					// report
					report()

					os.Exit(1)
				}
			}
		}()
	}

	// catch signal: interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		sig := <-sigChan
		fmt.Printf("quit by %s \n", sig.String())
		// report
		report()

		os.Exit(0)
	}()

	if <-quit {
	}
}

// print result to stdout
func report() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"name", "result"})

	data := [][]string{
		[]string{"limit request times", toString(limitRequestTimes)},
		[]string{"total request times", toString(totalRequestTimes)},
		[]string{"recover times", toString(recoverTimes)},
		[]string{"panic times", toString(panicTimes)},
		[]string{"max panic duration", toString(panics[len(panics)-1] - panics[0])},
		[]string{"concurrent", toString(co)},
	}

	average, levels := analysis()
	data = append(data, []string{"average", toString(average)})

	for k, v := range levels {
		data = append(data, []string{k, toString(v)})
	}

	for _, v := range data {
		table.Append(v)
	}

	fmt.Println()
	table.Render()
}

// parse args
func parse(args map[string]interface{}) {
	for k, v := range args {
		switch k {
		case "--headers":
			if v != nil {
				headers = v.(string)
			}
		case "--body":
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
		case "--limit":
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

func toString(value interface{}) string {
	var v string

	switch value.(type) {
	case string:
		v, _ = value.(string)
	case int:
		v = strconv.Itoa(value.(int))
	case int32:
		v = strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		v = strconv.FormatInt(value.(int64), 10)
	}

	return v
}
