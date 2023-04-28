package json

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewFetcher[T any](client *http.Client, baseURL string) *fetcher[T] {
	return &fetcher[T]{client, baseURL}
}

type fetcher[T any] struct {
	client  *http.Client
	baseURL string
}

func (f *fetcher[T]) FetchById(id string) (T, error) {
	var result T
	method := "GET"
	url := f.baseURL + "/" + id
	data, err := f.send(method, url)
	if err != nil {
		return result, fmt.Errorf("could not %+v %+v - %w", method, url, err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf(
			"could not unmarshal response body for %+v %+v - %w",
			method,
			url,
			err,
		)
	}
	return result, nil
}

func (f *fetcher[T]) FetchWhere(query string) ([]T, error) {
	var result []T
	method := "GET"
	url := f.baseURL + "?" + query
	data, err := f.send(method, url)
	if err != nil {
		return result, fmt.Errorf("could not %+v %+v - %w", method, url, err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf(
			"could not unmarshal response body for %+v %+v - %w",
			method,
			url,
			err,
		)
	}
	return result, err
}

func (f *fetcher[_]) send(method string, uri string) ([]byte, error) {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("could not construct request %w", err)
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get a response %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
