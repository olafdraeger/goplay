package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"lexemo/utilities/general"
	"net/http"
	"os"
	"strings"
)

//var baseUrl = []string{"https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:"}
var baseUrl = []string{"https://eur-lex.europa.eu/legal-content/EN/TXT/HTML/?uri=CELEX:"}
//var baseUrl = []string{"http://publications.europa.eu/resource/celex/"}

func main() {
	for _, eurlexId := range os.Args[1:] {
		url := strings.Join(append(baseUrl, eurlexId), "")
		fmt.Printf("Downloading: %v\nFrom: %v\n", eurlexId, url)
		resp, err := http.Get(url)
		generaltools.Check(err)
		defer resp.Body.Close()

		fmt.Printf("HTTP Status: %s\n", resp.Status)
		fmt.Printf("Download complete...\n\n")
		b, err := ioutil.ReadAll(resp.Body)
		bString := string(b)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		//doc, err := html.Parse(buf)
		//generaltools.Check(err)
		//var f func(*html.Node)
		//f = func(n *html.Node) {
		//	if n.Type == html.ElementNode && n.Data == "a" {
		//		for i, attr := range n.Attr {
		//			if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
		//				fmt.Printf("attribute index: %v\nattribute data: %v\n\n", i, attr)
		//			}
		//		}
		//	}
		//	for c := n.FirstChild; c != nil; c = c.NextSibling {
		//		f(c)
		//	}
		//}
		//f(doc)

		//z := html.NewTokenizer(buf)
		//for {
		//	tt := z.Next()
		//	if tt == html.ErrorToken {
		//		fmt.Println(tt.String())
		//		return
		//	}
		//	//fmt.Printf("current token: %v\n",tt.String())
		//	txtContent := strings.TrimSpace(html.UnescapeString(string(z.Text())))
		//
		//	fmt.Printf("tag: %v", tt.)
		//	if len(txtContent) > 0 {
		//		fmt.Printf("Content: %v\n", txtContent)
		//	}
		//}
		//resp.Body.Close()
		//generaltools.Check(err)
		//fmt.Printf("%s", b)
		//io.Copy(os.Stdout, resp.Body)

		var inBody = false
		domDocTest := html.NewTokenizer(strings.NewReader(bString))
		previousStartTokenTest := domDocTest.Token()
	loopDomTest:
		for {
			tt := domDocTest.Next()
			switch {
			case tt == html.ErrorToken:
				break loopDomTest // End of the document,  done
			case tt == html.StartTagToken:
				previousStartTokenTest = domDocTest.Token()
			case tt == html.TextToken:
				if previousStartTokenTest.Data == "body" {
					inBody = true
				}
				TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
				if inBody && len(TxtContent) > 0 {
					fmt.Printf("%s\t%s\n",domDocTest.Token().Data, TxtContent)
				}
			}
		}
	}
}
