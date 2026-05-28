package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL      string
	Status   int
	Bytes    int
	Duration time.Duration
	Err      error
}

var client = &http.Client{Timeout: 10 * time.Second}

func fetch(url string) Result {
	start := time.Now()
	resp, err := client.Get("https://" + url)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	defer resp.Body.Close()

	buf := make([]byte, 32*1024)
	total := 0
	for {
		n, err := resp.Body.Read(buf)
		total += n
		if err != nil {
			break
		}
	}
	return Result{
		URL:      url,
		Status:   resp.StatusCode,
		Bytes:    total,
		Duration: time.Since(start),
	}
}

func printResult(r Result) {
	if r.Err != nil {
		fmt.Printf("✗ %-25s ERR  %s\n", r.URL, r.Err)
		return
	}
	fmt.Printf("✓ %-25s %d  %.1fKB  %s\n",
		r.URL, r.Status, float64(r.Bytes)/1024, r.Duration)
}

func workerPool(urls []string, workers int) []Result {
	jobs := make(chan string, len(urls))
	results := make(chan Result, len(urls))

	// start N workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				results <- fetch(url)
			}
		}()
	}

	// send all jobs
	for _, url := range urls {
		jobs <- url
	}
	close(jobs) // signal workers: no more jobs

	// wait for all workers to finish, then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results
	all := []Result{}
	for r := range results {
		all = append(all, r)
	}
	return all
}

func main() {
	urls := []string{
		"go.dev",
		"example.com",
		"github.com",
		"golang.org",
		"httpbin.org/get",
		"cloudflare.com",
		"wikipedia.org",
		"stackoverflow.com",
		"thisurldoesnotexist123.xyz",
		"reddit.com",
	}

	const workers = 3
	fmt.Printf("Scraping %d URLs with %d workers...\n\n", len(urls), workers)

	start := time.Now()
	results := workerPool(urls, workers)

	ok, failed := 0, 0
	for _, r := range results {
		printResult(r)
		if r.Err != nil {
			failed++
		} else {
			ok++
		}
	}

	fmt.Printf("\nDone. %d ok, %d failed, total: %s\n",
		ok, failed, time.Since(start))
}