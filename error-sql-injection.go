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

// ANSI escape codes for text colors
const (
	redColor   = "\033[91m"
	greenColor = "\033[92m"
	resetColor = "\033[0m"
)

var (
	urlFile         string
	payloadFile     string
	wordToMatchFile string
	outputFile      string
	verbose         bool
	threads         int
)

func init() {
	flag.StringVar(&urlFile, "u", "", "File containing a list of URLs")
	flag.StringVar(&payloadFile, "p", "", "File containing payloads")
	flag.StringVar(&wordToMatchFile, "w", "", "File containing words to match")
	flag.StringVar(&outputFile, "o", "output.txt", "Output file to store results")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.IntVar(&threads, "t", 20, "Number of threads (concurrent requests)")
	flag.Parse()
}

func main() {
	// Define color constants for your banner
	CYAN := "\033[96m"
	NC := "\033[0m"

	// Add your banner here
	fmt.Print(CYAN, `
 
____ ____ ____ ____ ____    ____ ____ _       _ _  _  _ ____ ____ ___ _ ____ _  _ 
|___ |__/ |__/ |  | |__/ __ [__  |  | |    __ | |\ |  | |___ |     |  | |  | |\ | 
|___ |  \ |  \ |__| |  \    ___] |_\| |___    | | \| _| |___ |___  |  | |__| | \| 
                                                                                  
 
`, NC)

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

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threads)

	for _, url := range urls {
		for _, payload := range payloads {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(url, payload string) {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				fullURL := url + payload

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

				if vulnerable {
					color := redColor
					status := "Vulnerable"
					result := fmt.Sprintf("%s - %s%s%s\n", fullURL, color, status, resetColor)
					fmt.Print(result) // Print to console

					// Write only vulnerable URLs to the output file
					if err := writeToFile(fullURL, outputFile); err != nil {
						fmt.Printf("Error writing to output file: %v\n", err)
					}
				}
			}(url, payload)
		}
	}

	wg.Wait()
	close(semaphore)
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
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
