package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"lexemo/utilities/general"
	"net/http"
	"os"
	"strings"
)

var baseUrl = []string{"https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:"}

func main() {
	for _, eurlexId := range os.Args[1:] {
		url := strings.Join(append(baseUrl, eurlexId), "")
		fmt.Printf("Downloading: %v\nFrom: %v\n", eurlexId, url)
		resp, err := http.Get(url)
		generaltools.Check(err)
		fmt.Printf("Download complete...\n\n")
		//b, err := ioutil.ReadAll(resp.Body)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		doc, err := html.Parse(buf)
		generaltools.Check(err)
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for i, attr := range n.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
						fmt.Printf("attribute index: %v\nattribute data: %v\n\n", i, attr)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		//z := html.NewTokenizer(buf)
		//for {
		//	tt := z.Next()
		//	if tt == html.ErrorToken {
		//		fmt.Println(tt.String())
		//		return
		//	}
		//	fmt.Printf("current token: %v\n",tt.String())
		//}
		//resp.Body.Close()
		//generaltools.Check(err)
		//fmt.Printf("HTTP Status: %s\n", resp.Status)
		//fmt.Printf("%s", b)
		//io.Copy(os.Stdout, resp.Body)
	}
}
