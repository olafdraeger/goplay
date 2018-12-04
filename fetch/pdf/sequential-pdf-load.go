package main

import (
	"bufio"
	"fmt"
	"github.com/cavaliercoder/grab"
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
type WriteCounter1 struct {
	Total uint64
}

func (wc *WriteCounter1) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter1) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func pullOnePdf(url string, eurlexId string) {
	resp, err := http.Get(url)

	//defer func() {
	//	//Notify that we're done after this function
	//	chFinished <- true
	//}()

	if err != nil {
		fmt.Println("Error: Failed to download PDF: \"", url, "\"")
		return
	}

	fmt.Println(resp.Uncompressed)
	b := resp.Body
	defer b.Close()

	fname := eurlexId + "-content.pdf"
	pdfFile, err := os.Create(fname)
	generaltools.Check(err)
	defer pdfFile.Close()


	textFname := eurlexId + ".txt"
	textF, err := os.Create(textFname)
	generaltools.Check(err)
	defer textF.Close()

	//pdfFile, err := os.Open("CELEX32013L0036ENTXT.pdf")


	counter := &WriteCounter1{}
	//_, err = io.Copy(pdfFile, b)
	_, err1 := io.Copy(pdfFile, io.TeeReader(b, counter))
	generaltools.Check(err1)
	fmt.Print("\n")

	pdfReader, err := pdf.NewPdfReader(pdfFile)
	generaltools.Check(err)

	isEncrypted, err := pdfReader.IsEncrypted()
	generaltools.Check(err)

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		generaltools.Check(err)
	}

	numPages, err := pdfReader.GetNumPages()
	generaltools.Check(err)

	fmt.Printf("-----------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("-----------------------\n")
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

	//ch <- url
}

func savePDF(url string, eurlexId string) string {
	// create new pdf file
	pdfFile, err := os.Create(eurlexId + ".pdf")
	generaltools.Check(err)
	defer pdfFile.Close()

	// download content
	//resp, err := http.Get(url)
	//generaltools.Check(err)
	//defer resp.Body.Close()

	resp, err := grab.Get(".", url)
	generaltools.Check(err)

	fmt.Println("Saved pdf to: " + resp.Filename)

	return resp.Filename

	// write body to file
	//_, err = io.Copy(pdfFile, resp.Body)
	//generaltools.Check(err)
}

func extractText(pdfFileName string, celexId string) string {
	pdfFile, err := os.Open(pdfFileName)
	generaltools.Check(err)
	defer pdfFile.Close()

	textFile, err := os.Create(celexId + ".txt")
	generaltools.Check(err)
	defer textFile.Close()

	pdfReader, err := pdf.NewPdfReader(pdfFile)
	generaltools.Check(err)

	isEncrypted, err := pdfReader.IsEncrypted()
	generaltools.Check(err)

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		generaltools.Check(err)
	}

	numPages, err := pdfReader.GetNumPages()
	generaltools.Check(err)

	fmt.Printf("-----------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("-----------------------\n")
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
		textFile.WriteString(txt)

	}
	return textFile.Name()
}

func parseLegalText(fileName string) {
	f, err := os.Open(fileName)
	generaltools.Check(err)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(os.Stderr, "reading file: ", err)
	}
}

func main() {
	// https://eur-lex.europa.eu/legal-content/EN/TXT/PDF/?uri=CELEX:32014R1286&from=EN
	eurLexUrl := "https://eur-lex.europa.eu/legal-content/EN/TXT/PDF/?uri=CELEX:"
	eurLexPdfSuffix := "&from=EN"
	seedCelexIds := os.Args[1:]
	var seedPdfUrls []string

	for _, id := range seedCelexIds {
		seedPdfUrls = append(seedPdfUrls, eurLexUrl+id+eurLexPdfSuffix)
	}

	for i, url := range seedPdfUrls {
		//pullOnePdf(url, seedCelexIds[i])
		pdfFileName := savePDF(url, seedCelexIds[i])
		txtFileName := extractText(pdfFileName, seedCelexIds[i])
		parseLegalText(txtFileName)
	}

}
