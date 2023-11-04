package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    var (
        urlFile   string
        payloadFile string
        errorsFile string
        verbose   bool
        outputFile string
    )

    flag.StringVar(&urlFile, "u", "url.txt", "File containing the list of URLs")
    flag.StringVar(&payloadFile, "p", "payload.txt", "File containing the payload to append")
    flag.StringVar(&errorsFile, "e", "errors.txt", "File to store the summarized results")
    flag.BoolVar(&verbose, "v", false, "Enable vorvos output")
    flag.StringVar(&outputFile, "o", "output.txt", "File to store detailed results")
    flag.Parse()

    // Read the payload from payload.txt
    payload, err := ioutil.ReadFile(payloadFile)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Read URLs from url.txt
    urlList, err := readLines(urlFile)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Open the errors.txt file for writing summarized results
    errorsOutput, err := os.Create(errorsFile)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer errorsOutput.Close()

    // Open the output.txt file for writing detailed results
    outputFileHandle, err := os.Create(outputFile)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer outputFileHandle.Close()

    for _, url := range urlList {
        // Make a GET request to the URL
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println(err)
            continue
        }
        defer resp.Body.Close()

        // Get the initial content length
        initialContentLength := resp.ContentLength

        // Append the payload and make another GET request
        payloadURL := url + string(payload)
        resp, err = http.Get(payloadURL)
        if err != nil {
            fmt.Println(err)
            continue
        }
        defer resp.Body.Close()

        // Get the final content length
        finalContentLength := resp.ContentLength

        // Compare the content lengths
        if initialContentLength != finalContentLength {
            errorsOutput.WriteString("URL: " + url + " - Vulnerable\n")
            errorsOutput.WriteString("Initial Content Length: " + fmt.Sprintf("%d", initialContentLength) + "\n")
            errorsOutput.WriteString("Final Content Length: " + fmt.Sprintf("%d", finalContentLength) + "\n\n")
        } else {
            errorsOutput.WriteString("URL: " + url + " - Not Vulnerable\n\n")
        }

        // Write detailed results to the output file if verbose mode is enabled
        if verbose {
            detailedResult := "URL: " + url + "\nInitial Content Length: " + fmt.Sprintf("%d", initialContentLength) + "\nFinal Content Length: " + fmt.Sprintf("%d", finalContentLength) + "\nResponse Status Code: " + resp.Status + "\n\n"
            _, err := outputFileHandle.WriteString(detailedResult)
            if err != nil {
                fmt.Println(err)
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
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

