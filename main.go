package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	start := time.Now()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://example.com")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	stop := time.Now()
	elapsed := stop.Sub(start)

	fmt.Println("Response spend time:", elapsed)
	fmt.Println("Response status:", resp.Status)

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		fmt.Println("Server is up.")
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		fmt.Println("Client error.")
	case resp.StatusCode >= 500:
		fmt.Println("Server error.")
	default:
		fmt.Println("Another problem.")
	}

}
