package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// ANSI escape codes for text colors
const (
	redColor   = "\033[91m"
	greenColor = "\033[92m"
	resetColor = "\033[0m"
)

func main() {
	var urlFile, payloadFile, wordToMatchFile, outputFile string
	var verbose bool

	flag.StringVar(&urlFile, "u", "", "File containing a list of URLs")
	flag.StringVar(&payloadFile, "p", "", "File containing payloads")
	flag.StringVar(&wordToMatchFile, "w", "", "File containing words to match")
	flag.StringVar(&outputFile, "o", "output.txt", "Output file to store results")
	flag.BoolVar(&verbose, "v", false, "Verbose output")

	flag.Parse()

	if urlFile == "" || payloadFile == "" || wordToMatchFile == "" {
		fmt.Println("Please provide all required input files.")
		return
	}

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

	for _, url := range urls {
		for _, payload := range payloads {
			fullURL := url + payload

			body, err := fetchURL(fullURL)
			if err != nil {
				fmt.Printf("Error fetching URL %s: %v\n", fullURL, err)
				continue
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
