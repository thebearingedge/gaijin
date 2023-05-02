package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func SendRequest[T any](client *http.Client, url string) (T, error) {
	var val T
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return val, fmt.Errorf("instantiating request %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return val, fmt.Errorf("sending request %w", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return val, fmt.Errorf("reading body %w", err)
	}
	if err := json.Unmarshal(data, &val); err != nil {
		return val, fmt.Errorf("deserializing json %w", err)
	}
	return val, err
}

type Result[T any] struct {
	Ok  T
	Err error
}

func SendRequestWithChan[T any](client *http.Client, url string, ch chan<- Result[T]) {
	defer close(ch)
	val, err := SendRequest[T](client, url)
	ch <- Result[T]{val, err}
}

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func concurrently[T any](urls []string) string {
	ch1 := make(chan Result[T])
	ch2 := make(chan Result[T])
	go SendRequestWithChan(http.DefaultClient, urls[0], ch1)
	go SendRequestWithChan(http.DefaultClient, urls[1], ch2)
	fst, snd := <-ch1, <-ch2
	return fmt.Sprintf("first result: %+v\n\n second result: %+v\n\n", fst, snd)
}

func sequentially[T any](urls []string) string {
	p1, err := SendRequest[T](http.DefaultClient, urls[0])
	fst := Result[T]{p1, err}
	p2, err := SendRequest[T](http.DefaultClient, urls[1])
	snd := Result[T]{p2, err}
	return fmt.Sprintf("first result: %+v\n\n second result: %+v\n\n", fst, snd)
}

func main() {
	start := time.Now().UnixMilli()
	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
	}

	// report := sequentially[Post](urls)
	report := concurrently[Post](urls)

	fmt.Println(report)
	fmt.Printf("completed in %v milliseconds\n", time.Now().UnixMilli()-start)
}
