package main

import (
	"fmt"
	"io"
	"lexemo/utilities/general"
	"net/http"
	"os"
	"strings"
	"time"
)

var eurLexUrl = []string{"https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:"}
var dirPath = "./html/"

func main() {
	if exists, _ := generaltools.Exists(dirPath); !exists {
		err := os.MkdirAll(dirPath, 0755)
		generaltools.Check(err)
	}
	start := time.Now()
	ch := make(chan string)
	for _, eurlexId := range os.Args[1:] {
		url := strings.Join(append(eurLexUrl, eurlexId), "")
		go fetchAll(url, eurlexId, ch)
	}
	for range os.Args[1:] {
		fmt.Println((<-ch))
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// fetches all eurlex articles by Id's as provided via command line
func fetchAll(url string, eurlexId string, ch chan<- string) {

	filePath := dirPath + eurlexId + ".html"

	// if html file exists, delete it
	if exists, _ := generaltools.Exists(filePath); exists {
		err := os.Remove(filePath)
		generaltools.Check(err)
	}

	//create new html file
	f, err := os.Create(filePath)
	generaltools.Check(err)
	defer f.Close()

	// start time keeping to check download time and retrieve content
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	nbytes, err := io.Copy(f, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)

}
