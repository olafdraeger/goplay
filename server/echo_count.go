package main

import (
	"fmt"
	"github.com/traefik/log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex
var count int

func main() {

	lissa := func(w http.ResponseWriter, r *http.Request) {
		var cycles float64 = 5
		if r.URL.RawQuery != "" {
			str := strings.Split(r.URL.RawQuery, "=")
			c, err := strconv.ParseFloat(str[1], 64)
			if err != nil {
				fmt.Printf("No Prameter provided")
			}
			cycles = c
		}
		lissajous(w, cycles)
	}

	http.HandleFunc("/", lissa)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/recho", echoRequest)
	http.HandleFunc("/lissa", lissa)
	//http.HandleFunc("/lissa", func(w http.ResponseWriter, r *http.Request) {
	//	lissajous(w)
	//})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler1(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()

	fmt.Fprintf(w, "URL.Path = %q\nParam: %q\n", r.URL.Path, r.URL.RawQuery)
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Counter: %d\n", count)
	fmt.Printf("Counter: %d\n", count)
	mu.Unlock()
}

func echoRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q = %q\n", k, v)
	}
}
