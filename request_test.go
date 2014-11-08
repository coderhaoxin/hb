package main

import "testing"

func TestRequest(t *testing.T) {
	request("GET", "http://www.baidu.com", "", "")
}
