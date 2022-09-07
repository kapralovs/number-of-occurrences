package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type UrlInfo struct {
	Url   string
	Count int
}

func main() {
	var (
		total      int
		wg         sync.WaitGroup
		listOfUrls = []string{
			"https://go.dev",
			"https://gobyexample.com",
			"http://golang-book.ru",
		}
	)
	ch := make(chan UrlInfo, len(listOfUrls))

	for _, url := range listOfUrls {
		body, err := makeRequest(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		info := UrlInfo{
			Url: url,
		}
		wg.Add(1)
		go findAllOccurrences(info, body, ch, &wg)
	}

	wg.Wait()
	close(ch)

	for {
		if info, opened := <-ch; opened {
			fmt.Printf("Count for %s: %d\n", info.Url, info.Count)
			total += info.Count
		} else {
			break
		}
	}
	fmt.Printf("Total: %d\n", total)
}

func makeRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func findAllOccurrences(info UrlInfo, body string, ch chan UrlInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	count := strings.Count(body, "Go")
	info.Count = count
	ch <- info
}
