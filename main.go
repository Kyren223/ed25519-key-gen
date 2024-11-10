package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

const goroutines = 12

func main() {
	keywordsStrings := getKeywords()
	var keywords [][]byte
	for _, keyword := range keywordsStrings {
		keywords = append(keywords, []byte(keyword))
	}
	if len(keywords) == 0 {
		fmt.Println("No keywords")
		return
	}

	output, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Println("File error:", err)
	}

	fmt.Println("Keywords: ", keywordsStrings)
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
					sshPubKey, err := ssh.NewPublicKey(pubKey)
					if err != nil {
						fmt.Println("Encountered error:", err)
						break
					}
					pub := ssh.MarshalAuthorizedKey(sshPubKey)[37:]
					pub = pub[:len(pub)-1]
					generated[i]++
					for _, keyword := range keywords {
						if bytes.Contains(pub, keyword) {
							mu.Lock()
							found++
							privBase64 := base64.StdEncoding.EncodeToString(privKey)
							pubBase64 := base64.StdEncoding.EncodeToString(pubKey)
							fmt.Printf("%X Found %s %s\n", i, string(keyword), string(pub))
							_, err := output.WriteString(fmt.Sprintln(string(keyword), "AAAAC3NzaC1lZDI1NTE5AAAAI"+string(pub), privBase64, pubBase64))
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

// Failed, no keys still contained their keywords
func update() {
	a, _ := os.ReadFile("a.txt")
	lines := strings.Split(strings.ReplaceAll(string(a), "\r\n", "\n"), "\n")
	var newLines []string
	for _, line := range lines {
		split := strings.Split(line, " ")
		if len(split) != 3 {
			continue
		}
		keyword, pub, _ := split[0], split[1], split[2]
		b, _ := base64.StdEncoding.DecodeString(pub)
		key, _ := ssh.NewPublicKey(ed25519.PublicKey(b))
		sshKey := ssh.MarshalAuthorizedKey(key)
		if bytes.Contains(sshKey, []byte(keyword)) {
			newLines = append(newLines, string(sshKey))
		}
	}
	os.WriteFile("b.txt", []byte(strings.Join(newLines, "\n")), 0o644)
}
