package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

const goroutines = 12

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

	fmt.Println("Keywords: ", keywords)
	fmt.Println("Goroutines: ", goroutines)
	fmt.Println("Starting to search...")

	var mu sync.Mutex
	var wg sync.WaitGroup
	var generated [goroutines]int64
	found := 0
	start := time.Now()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt)

			for {
				select {
				case <-signalChan:
					return
				default:
					pubKey, privKey, err := ed25519.GenerateKey(nil)
					if err != nil {
						fmt.Println("Encountered error:", err)
						break
					}
					generated[i]++
					pubBase64 := base64.StdEncoding.EncodeToString(pubKey)
					for _, keyword := range keywords {
						if strings.Contains(pubBase64, keyword) {
							mu.Lock()
							found++
							privBase64 := base64.StdEncoding.EncodeToString(privKey)
							fmt.Println(i, "Found", keyword, pubBase64)
							_, err := output.WriteString(fmt.Sprintln(keyword, pubBase64, privBase64))
							if err != nil {
								fmt.Println("Error writing string to output file:", err)
								fmt.Println("Private Key:", privKey)
							}
							mu.Unlock()
							break
						}
					}
				}
			}
		}()
	}

	wg.Wait()

	searched := int64(0)
	for _, gen := range generated {
		searched += gen
	}

	elapsed := time.Since(start)
	fmt.Println("Completed Search")
	fmt.Println("Time Elapsed:", elapsed)
	fmt.Println("Searched:", searched, "Found:", found)
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
