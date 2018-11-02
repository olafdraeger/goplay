package main

import (
	"fmt"
	"io"
	"lexemo/utilities/general"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http://") {
			elems := []string{"http://", url}
			url = strings.Join(elems, "")
		}
		resp, err := http.Get(url)
		generaltools.Check(err)
		//b, err := ioutil.ReadAll(resp.Body)
		//resp.Body.Close()
		//generaltools.Check(err)
		fmt.Printf("HTTP Status: %s\n", resp.Status)
		//fmt.Printf("%s", b)
		io.Copy(os.Stdout, resp.Body)
	}
}
