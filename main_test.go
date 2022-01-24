package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestCallURL(t *testing.T) {
	t.Run("Sends error with invalid URL", func(t *testing.T) {
		limit := 1
		count := 1
		testURL := "TEST**"

		errc := make(chan error, limit)
		resc := make(chan result, limit)
		URLc := make(chan string, count)
		wg := sync.WaitGroup{}
		wg.Add(1)
		h := &helper{
			errc: errc,
			resc: resc,
			URLc: URLc,
			wg:   wg,
		}

		client := http.Client{}
		callURL(testURL, h, &client)
		err := <-h.errc
		if err == nil {
			t.Fail()
		}
	})
	t.Run("Sends error with connection refused", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		ts.Close()

		limit := 1
		count := 1
		testURL := ts.URL

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
		h.wg.Add(1)

		client := ts.Client()
		callURL(testURL, h, client)
		<-h.errc
	})

	t.Run("Decrements wait group only once", func(t *testing.T) {
		defer func() {
			err := recover()
			if err != "sync: negative WaitGroup counter" {
				t.Fatalf("Unexpected panic: %#v", err)
			}
		}()

		limit := 1
		count := 1
		testURL := "TEST***"

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
		h.wg.Add(1)
		callURL(testURL, h, &http.Client{})
		wg.Done()
		t.Fatal("Should panic")
	})
}

func BenchmarkProcessAll(b *testing.B) {
	URLs := []string{"google.com", "facebook.com", "reddit.com", "youtube.com", "twitter.com"}
	b.Run("Benchmark limit = 1", func(b *testing.B) {
		limit := 1
		processAll(limit, URLs)
	})
	b.Run("Benchmark full capacity", func(b *testing.B) {
		limit := len(URLs)
		processAll(limit, URLs)
	})
}
