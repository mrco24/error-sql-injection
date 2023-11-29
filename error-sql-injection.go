package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// ANSI escape codes for text colors
const (
	redColor    = "\033[91m"
	greenColor  = "\033[92m"
	yellowColor = "\033[93m"
	resetColor  = "\033[0m"
)

var (
	url        string
	urlFile    string
	outputFile string
	verbose    bool
	threads    int
)

// Default payloads
var defaultPayloads = []string{
	"%27%22%60",
	"'",
}

func init() {
	flag.StringVar(&url, "u", "", "Single target URL")
	flag.StringVar(&urlFile, "f", "", "File containing a list of URLs")
	flag.StringVar(&outputFile, "o", "output.txt", "Output file to store results")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.IntVar(&threads, "t", 20, "Number of threads (concurrent requests)")
	flag.Parse()
}

func main() {
	var urls []string

	if url != "" {
		urls = append(urls, url)
	} else if urlFile != "" {
		urlsFromFile, err := readLines(urlFile)
		if err != nil {
			fmt.Printf("Error reading URLs from %s: %v\n", urlFile, err)
			return
		}
		urls = append(urls, urlsFromFile...)
	}

	for _, u := range urls {
		for _, payload := range defaultPayloads {
			fullURL := u + payload

			body, err := fetchURL(fullURL)
			// Print the request URL prefix with color, regardless of success or failure
			fmt.Printf("%sRequest URL:%s %s\n", yellowColor, resetColor, fullURL)

			if err != nil {
				fmt.Printf("Error fetching URL %s: %v\n", fullURL, err)
			}

			vulnerable := false
			wordsToMatch := []string{
				"mysql",
				"fetch_array",
				"SQL syntax",
				"500 Internal Server Error",
				"mysqli",
				"Access Database Engine",
				"SQLite",
				"Sybase",
				"Server message",
				"Oracle",
				"Driver",
			}

			for _, word := range wordsToMatch {
				if strings.Contains(body, word) {
					vulnerable = true
					break
				}
			}

			var color, status string
			if vulnerable {
				color = redColor
				status = "Vulnerable"
			} else {
				color = greenColor
				status = "Next"
			}

			// Print the full URL and status with color
			result := fmt.Sprintf("Full URL: %s - %s%s%s\n", fullURL, color, status, resetColor)
			fmt.Print(result) // Print to console

			// Write only vulnerable URLs to the output file
			if vulnerable {
				if err := writeToFile(fullURL+" - "+status, outputFile); err != nil {
					fmt.Printf("Error writing to output file: %v\n", err)
				}
			}
		}
	}
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	lines = strings.Split(string(content), "\n")

	return lines, nil
}

func fetchURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func writeToFile(output string, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(output + "\n"); err != nil {
		return err
	}

	return nil
}
