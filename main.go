package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type result struct {
	url  string
	hash hash.Hash
}
type helper struct {
	resc chan result
	errc chan error
	URLc chan string
	wg   sync.WaitGroup
}

func callURL(URL string, h *helper) {
	defer h.wg.Done()
	client := http.Client{}
	res := result{}
	reqURL, err := url.Parse(URL)
	if err != nil {
		h.errc <- fmt.Errorf("error parsing request URL %s:  %w", URL, err)
		return
	}
	scheme := reqURL.Scheme
	if scheme != "http" && scheme != "https" {
		reqURL.Scheme = "https"
	}
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		h.errc <- fmt.Errorf("error making request URL %s:  %w", URL, err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		h.errc <- fmt.Errorf("error doing request URL %s:  %w", URL, err)
		return
	}
	hash := md5.New()
	_, err = io.Copy(hash, resp.Body)
	if err != nil {
		h.errc <- fmt.Errorf("error generating hash URL %s:  %w", URL, err)
		return
	}
	res.url = reqURL.String()
	res.hash = hash
	h.resc <- res
}

func worker(h *helper) {
	for URL := range h.URLc {
		callURL(URL, h)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func processAll(limit int, URLs []string) {
	count := len(URLs)
	workersCount := min(limit, count)
	fmt.Printf("Limit = %d \n", limit)
	fmt.Printf("URL count = %d \n", count)
	fmt.Printf("Workers %d \n", workersCount)
	errc := make(chan error, limit)
	resc := make(chan result, limit)
	URLc := make(chan string, count)
	wg := sync.WaitGroup{}
	h := &helper{
		errc: errc,
		resc: resc,
		URLc: URLc,
		wg:   wg,
	}

	for i := 0; i < workersCount; i++ {
		go worker(h)
	}

	for _, URL := range URLs {
		h.wg.Add(1)
		h.URLc <- URL
	}
	close(h.URLc)

	for i := 0; i < count; i++ {
		select {
		case err := <-h.errc:
			if err != nil {
				fmt.Printf("Error %s \n", err.Error())
			}
		case res := <-h.resc:
			fmt.Printf("%s %x \n", res.url, res.hash.Sum(nil))
		}
	}
	h.wg.Wait()
	close(h.errc)
	close(h.resc)
}

func main() {
	// parse the parallel flag
	limit := flag.Int("parallel", 10, "limit of parallel requests")
	flag.Parse()
	processAll(*limit, flag.Args())
}
