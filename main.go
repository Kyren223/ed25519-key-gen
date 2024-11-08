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
	keywords := []string{
		"Kyren",
		"kyren",
		"KYREN",
		"Kyren223",
		"kyren223",
		"KYREN223",
		"Banana",
		"banana",
	}

	output, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Println("File error:", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

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
