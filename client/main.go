package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
)

const (
	fmtURL = "https://xkcd.com/%d/info.0.json"
)

// FetchAll Launches parallel goroutines to fetch data from individual pages.
// It doesn't know apriori how many pages it will fetch. It starts at 1 and increments.
// Stops at the first 404 response that it gets.
//
// Deals with error internally by printing them out.
// TODO maybe implement a better error handling stategy (like reporting them to the caller)
func FetchAll(out io.Writer, workers int) {

	ctx, cancel := context.WithCancel(context.Background())

	ids := make(chan int)
	data := make(chan *Result)
	errs := make(chan error)

	// start workers
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(ctx, &wg, ids, data, errs)
	}

	//wait for errors
	go func() {
		defer cancel()
		handleErrors(ctx, errs)
	}()

	// send ID's on the ids channel
	go createWorkload(ctx, ids)

	// do something with the data
	writeCtx, writeCtxCancel := context.WithCancel(context.Background())
	go func() {
		defer writeCtxCancel()

		s := make([]*Result, 0)
		for d := range data {
			s = append(s, d)
		}

		sortAscending(s)

		if err := json.NewEncoder(out).Encode(s); err != nil {
			log.Println(err)
		}
	}()

	// wait for workers to finish
	wg.Wait()

	// close the data channel and trigger outputting the data
	close(data)

	// wait for the data to be output
	<-writeCtx.Done()
}

func handleErrors(ctx context.Context, errs <-chan error) {
	for {
		select {
		case err := <-errs:
			// print the error
			log.Println(err)
			return
		case <-ctx.Done():
			return
		}
	}

}

func createWorkload(ctx context.Context, ids chan<- int) {
	counter := 0

	defer func() {
		close(ids)
	}()

	for {
		counter++
		if counter == 404 {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case ids <- counter:
		}
	}

}

func worker(ctx context.Context, wg *sync.WaitGroup, ids <-chan int, data chan<- *Result, errs chan error) {
	defer wg.Done()

	for id := range ids {
		r, err := query(id)
		if err != nil {
			select {
			case <-ctx.Done():
			case errs <- err:
			}
			continue
		}

		select {
		case <-ctx.Done():
			return
		case data <- r:
		}
	}
}

func query(id int) (*Result, error) {

	resp, err := http.Get(fmt.Sprintf(fmtURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code %d for id %d", resp.StatusCode, id)
	}

	var r Result
	err = json.NewDecoder(resp.Body).Decode(&r)

	return &r, err
}

func sortAscending(s []*Result) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Num < s[j].Num
	})
}
