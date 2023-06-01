package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	str := flag.String("str", "", "Input string without spaces")
	urlStr := flag.String("url", "", "URL to send the request")
	flag.Parse()

	if *str == "" || *urlStr == "" {
		log.Fatal("Both -str and -url flags are required")
	}

	resp, err := http.PostForm(*urlStr+"/api/substring", url.Values{"str": {*str}})
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Request failed with status code:", resp.StatusCode)
	}

	var result string
	_, err = fmt.Fscanf(resp.Body, "%s", &result)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	fmt.Println("Longest substring without repeating characters:", result)
}
