package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	keywords := getKeywords()
	if len(keywords) == 0 {
		fmt.Println("No keywords")
		return
	}

	output, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Println("File error:", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("Keywords: ", keywords)
	fmt.Println("Starting to search...")

	start := time.Now()
	generated := 0
	found := 0

outer:
	for {
		select {
		case <-signalChan:
			fmt.Println("Interrupt")
			break outer
		default:
			pubKey, privKey, err := ed25519.GenerateKey(nil)
			if err != nil {
				fmt.Println("Encountered error:", err)
				break
			}
			generated++
			pubBase64 := base64.StdEncoding.EncodeToString(pubKey)
			for _, keyword := range keywords {
				if strings.Contains(pubBase64, keyword) {
					found++
					privBase64 := base64.StdEncoding.EncodeToString(privKey)
					fmt.Println("Found", keyword, pubBase64)
					_, err := output.WriteString(fmt.Sprintln(keyword, pubBase64, privBase64))
					if err != nil {
						fmt.Println("Error writing string to output file:", err)
						fmt.Println("Private Key:", privKey)
					}
					break
				}
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Completed Search")
	fmt.Println("Time Elapsed:", elapsed)
	fmt.Println("Searched:", generated, "Found:", found)
}

func getKeywords() []string {
	var keywords []string
	args := os.Args[1:]
	if len(args) == 0 {
		file, err := os.ReadFile("input.txt")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		contents := strings.ReplaceAll(string(file), "\r\n", "\n")
		for _, keyword := range strings.Split(contents, "\n") {
			keyword = strings.TrimSpace(keyword)
			if len(keyword) != 0 {
				keywords = append(keywords, keyword)
			}
		}
	} else {
		for _, keyword := range args {
			keyword = strings.TrimSpace(keyword)
			if len(keyword) != 0 {
				keywords = append(keywords, keyword)
			}
		}
	}

	return keywords
}
