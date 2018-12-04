package main

import (
	"fmt"
	"golang.org/x/net/html"
	"lexemo/utilities/general"
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

func checkIfContent(t html.Token) (ok bool, div string) {
	// Iterating over token's attributes to find a div with ID £PP4Contents
	for _, d := range t.Attr {
		if d.Key == "id" && d.Val == "document1"{
			div = d.Val
			ok = true
		}
	}
	return
}

func crawlLinks(url string, eurlexId string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		//Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("Errpr: Failed to crawlLinks: \"", url, "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	fname := eurlexId + "-links.txt"
	f, err := os.Create(fname)
	generaltools.Check(err)
	defer f.Close()

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
				f.WriteString(url + "\n")
				ch <- url
			}


		}
	}
}

func crawlContent(url string, eurlexId string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		//Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("Error: Failed to crawl: \"", url, "\"")
		return
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	fname := eurlexId + "-content.txt"
	f, err := os.Create(fname)
	generaltools.Check(err)
	defer f.Close()

	isContent := false
	//domDocTest := html.NewTokenizer(strings.NewReader(s))
	previousStartTokenTest := z.Token()
loopDomTest:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // End of the document,  done
		case tt == html.StartTagToken:
			previousStartTokenTest = z.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" {
				continue
			}

			// once we found the PP4Content div id we reached the content section of relevance
			ok, _ := checkIfContent(previousStartTokenTest)
			if ok {
				isContent = true
			}

			if isContent {
				TxtContent := strings.TrimSpace(html.UnescapeString(string(z.Text())))
				if len(TxtContent) > 0 {
					for _, attr := range previousStartTokenTest.Attr {
						f.WriteString("Key: " + attr.Key + " - Value: " + attr.Val + "\n")
					}
					f.WriteString("Text: " + TxtContent + "\n")
					f.WriteString("=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+\n")
				}
			}
		}
	}
	ch <- url
}

//func pullPDF(url string, urlSuffix string, eurlexId string, ch chan string, chFinished chan bool) {
//	resp, err := http.Get(url)
//
//	defer func() {
//		//Notify that we're done after this function
//		chFinished <- true
//	}()
//
//	if err != nil {
//		fmt.Println("Error: Failed to download PDF: \"", url, "\"")
//		return
//	}
//
//	b := resp.Body
//	defer b.Close()
//
//	fname := eurlexId + "-content.pdf"
//	f, err := os.Create(fname)
//	generaltools.Check(err)
//	defer f.Close()
//}

func main() {
	eurLexUrl := "https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:"
	//eurLexPdfSuffix := "&from=EN"
	foundUrls := make(map[string]bool)
	seedCelexIds := os.Args[1:]
	var seedUrls[]string
	//var  seedPdfUrls []string

	for _, id := range seedCelexIds {
		seedUrls = append(seedUrls, eurLexUrl+id)
		//seedPdfUrls = append(seedPdfUrls, eurLexUrl+id+eurLexPdfSuffix)
	}

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// fire off as many goroutines as we have URLS to check
	for i, url := range seedUrls {
		//go crawlLinks(url, seedCelexIds[i], chUrls, chFinished)
		go crawlContent(url, seedCelexIds[i], chUrls, chFinished)
		//go pullPDF(url, eurLexPdfSuffix, seedCelexIds[i], chUrls, chFinished)
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
