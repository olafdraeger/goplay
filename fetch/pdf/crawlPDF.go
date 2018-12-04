package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	pdf "github.com/unidoc/unidoc/pdf/model"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	"strings"

	//"github.com/rsc/pdf"
	"lexemo/utilities/general"
	"net/http"
	"os"
)
//32015R0227  32009L0138 02009L0065-20140917 32011L0061 Priips: 32014R1286
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func pullPDF(url string, eurlexId string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		//Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("Error: Failed to download PDF: \"", url, "\"")
		return
	}

	fmt.Println(resp.Uncompressed)
	b := resp.Body
	defer b.Close()

	fname := eurlexId + "-content.pdf"
	f, err := os.Create(fname)
	generaltools.Check(err)
	defer f.Close()


	textFname := eurlexId + ".txt"
	textF, err := os.Create(textFname)
	generaltools.Check(err)
	defer textF.Close()

	//pdfFile, err := os.Open("CELEX32013L0036ENTXT.pdf")


	counter := &WriteCounter{}
	//_, err = io.Copy(f, b)
	_, err1 := io.Copy(textF, io.TeeReader(b, counter))
	generaltools.Check(err1)
	fmt.Print("\n")

	pdfReader, err := pdf.NewPdfReader(f)
	generaltools.Check(err)

	isEncrypted, err := pdfReader.IsEncrypted()
	generaltools.Check(err)

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		generaltools.Check(err)
	}

	numPages, err := pdfReader.GetNumPages()
	generaltools.Check(err)

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		generaltools.Check(err)

		contentStreams, err := page.GetContentStreams()
		generaltools.Check(err)

		// If the value is an array, the effect shall be as if all of the streams in the array were concatenated,
		// in order, to form a single stream.
		pageContentStr := ""
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}

		fmt.Printf("Page %d - content streams %d:\n", pageNum, len(contentStreams))
		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		txt, err := cstreamParser.ExtractText()
		generaltools.Check(err)
		textF.WriteString(txt)
		//fmt.Printf("\"%s\"\n", txt)
	}

	ch <- url
}

func main() {
	eurLexUrl := "https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:"
	eurLexPdfSuffix := "&from=EN"
	foundUrls := make(map[string]bool)
	seedCelexIds := os.Args[1:]
	var seedPdfUrls []string

	for _, id := range seedCelexIds {
		seedPdfUrls = append(seedPdfUrls, eurLexUrl+id+eurLexPdfSuffix)
	}

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// fire off as many goroutines as we have URLS to check
	for i, url := range seedPdfUrls {
		go pullPDF(url, seedCelexIds[i], chUrls, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedPdfUrls); {
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
