package main

import (
	"flag"
	"fmt"
)

func main() {
	// parse the parallel flag
	limit := flag.Int("parallel", 10, "limit of parallel requests")
	flag.Parse()
	fmt.Printf("Limit = %d \n", *limit)

	inputURLs := flag.Args()
	for _, url := range inputURLs {
		fmt.Printf("URLS = %s \n", url)
	}
}
