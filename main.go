package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	"strings"
)

func main() {
	// Print the banner
	printBanner()

	// Define command-line flags
	url := flag.String("u", "", "URL to test")
	headers := flag.String("H", "", "HTTP headers")
	successCode := flag.Int("success", http.StatusOK, "Expected HTTP success code")
	errorCode := flag.Int("error", http.StatusTooManyRequests, "Expected HTTP error code")
	duration := flag.Int("t", 0, "Duration to send requests")

	// Override the default usage message
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage:\n")
        fmt.Fprintf(os.Stderr, "  -u string\tURL to test\n")
        fmt.Fprintf(os.Stderr, "  -H string\tSpecify headers (e.g., \"Myheader: test\")\n")
        fmt.Fprintf(os.Stderr, "  --success int\tExpected HTTP success code (default 200)\n")
        fmt.Fprintf(os.Stderr, "  --error int\tExpected HTTP error code (default 429)\n")
        fmt.Fprintf(os.Stderr, "  -t int\tSet duration(seconds) to send maximum number of requests\n")
    }

    // Parse command-line flags
    flag.Parse()

    // If the -h flag is used, print the usage and exit
    if flag.Lookup("h") != nil {
        flag.Usage()
        os.Exit(0)
    }

	// Check if URL is provided
	if *url == "" {
		fmt.Println("Please provide a URL using the -u flag")
		os.Exit(1)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		os.Exit(1)
	}

	// Set HTTP headers if provided
	if *headers != "" {
		headerSlice := parseHeaders(*headers)
		for key, value := range headerSlice {
			req.Header.Set(key, value)		
		}
	}

	// Send the HTTP request for the specified duration
	endTime := time.Now().Add(time.Duration(*duration) * time.Second)
	counter := 0
	for time.Now().Before(endTime) {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to send HTTP request:", err)
			os.Exit(1)
		}

		counter++
		fmt.Print(time.Now().Format("2006-01-02 15:04:05"), " ")

		// Check the HTTP response code
		if resp.StatusCode == *successCode {
			fmt.Println("Success! Expected HTTP success code received:", resp.StatusCode)
		} else if resp.StatusCode == *errorCode {
			fmt.Println("Error! Expected HTTP error code received:", resp.StatusCode)
			fmt.Println("Rate limit hit: ",counter-1, " requests")
			os.Exit(1)
		} else {
			fmt.Println("Unexpected HTTP response code received:", resp.StatusCode)
			fmt.Print("\n\033[1A\033[K")
		}

		// Close the response body immediately after checking the status code
		resp.Body.Close()
	}
}

// Helper function to parse HTTP headers
func parseHeaders(headers string) map[string]string {
	headerSlice := make(map[string]string)
	headerList := strings.Split(headers, ",")

	for _, header := range headerList {
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) == 2 {
			key := strings.TrimSpace(headerParts[0])
			value := strings.TrimSpace(headerParts[1])
			headerSlice[key] = value
		}
	}

	return headerSlice
}

func printBanner() {
	redColor := "\033[91m"
	resetColor := "\033[0m"

	ascii := `
         ________             .____    .__        .__  __   
        /  _____/  ___________|    |   |__| _____ |__|/  |_ 
       /   \  ___ /  _ \_  __ \    |   |  |/     \|  \   __\
       \    \_\  (  <_> )  | \/    |___|  |  Y Y  \  ||  |  
        \______  /\____/|__|  |_______ \__|__|_|  /__||__|  
               \/                     \/        \/                                                                 
                                                             c=====e
                                                            	H
                                                            _,,_H__
      (__((__((___()                                       //|     |
     (__((__((___()()_____________________________________// |ACME |
    (__((__((___()()()------------------------------------'  |_____|
	
	`
	
	fmt.Println(redColor,ascii,resetColor)
}