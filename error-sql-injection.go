package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ... (constants and functions remain unchanged)

func main() {
	var urlFile, payloadFile, wordToMatchFile, outputFile string
	var verbose bool
	var threads int

	flag.StringVar(&urlFile, "u", "", "File containing a list of URLs")
	flag.StringVar(&payloadFile, "p", "", "File containing payloads")
	flag.StringVar(&wordToMatchFile, "w", "", "File containing words to match")
	flag.StringVar(&outputFile, "o", "output.txt", "Output file to store results")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.IntVar(&threads, "h", 1, "Number of threads (concurrent requests)")

	flag.Parse()

	if urlFile == "" || payloadFile == "" || wordToMatchFile == "" {
		fmt.Println("Please provide all required input files.")
		return
	}

	// ... (banner and reading input files remain unchanged)

	var wg sync.WaitGroup
	urls, err := readLines(urlFile)
	if err != nil {
		fmt.Printf("Error reading URLs from %s: %v\n", urlFile, err)
		return
	}

	payloads, err := readLines(payloadFile)
	if err != nil {
		fmt.Printf("Error reading payloads from %s: %v\n", payloadFile, err)
		return
	}

	wordsToMatch, err := readLines(wordToMatchFile)
	if err != nil {
		fmt.Printf("Error reading words to match from %s: %v\n", wordToMatchFile, err)
		return
	}

	for _, url := range urls {
		for _, payload := range payloads {
			fullURL := url + payload
			wg.Add(1)
			go func(url, fullURL, payload string, wordsToMatch []string) {
				defer wg.Done()
				body, err := fetchURL(fullURL)
				if err != nil {
					fmt.Printf("Error fetching URL %s: %v\n", fullURL, err)
					return
				}

				vulnerable := false
				for _, word := range wordsToMatch {
					if strings.Contains(body, word) {
						vulnerable = true
						break
					}
				}

				color := greenColor
				status := "Not Vulnerable"
				if vulnerable {
					color = redColor
					status = "Vulnerable"
				}
				fmt.Printf("%s - %s%s%s\n", fullURL, color, status, resetColor)
			}(url, fullURL, payload, wordsToMatch)
		}
	}

	wg.Wait()
}

// ... (readLines and fetchURL functions remain unchanged)
