package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

func getHref(t html.Token) (ok bool, href string) {
	// Iterating over token's attributes to find a href
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	// "bare" return will return the variables (ok, href) as defined in the
	// function definition
	return
}

func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		//Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("Errpr: Failed to crawl: \"", url, "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// end of document - done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// check what token type

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// extract href value is exists
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure we're looking at complete URLs
			hasProto := strings.HasPrefix(url, "http")
			if hasProto {
				ch <- url
			}

		}
	}
}

func main() {
	eurLexUrl := "https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:"
	foundUrls := make(map[string]bool)
	seedCelexIds := os.Args[1:]
	var seedUrls []string

	for _, id := range seedCelexIds {
		seedUrls = append(seedUrls, eurLexUrl+id)
	}

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// fire off as many goroutines as we have URLS to check
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	// Done parsing, print results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - ", url)
	}

	close(chUrls)
}
