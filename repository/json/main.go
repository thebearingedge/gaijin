package json

import (
	"encoding/json"
	"io"
	"net/http"
)

func NewFetcher[T any](client http.Client, baseURL string) *fetcher[T] {
	return &fetcher[T]{baseURL, client}
}

type fetcher[T any] struct {
	baseURL string
	client  http.Client
}

func (f fetcher[T]) FetchById(id string) (T, error) {
	var result T
	data, err := f.send("GET", f.baseURL+"/"+id)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (f fetcher[T]) FetchWhere(query string) ([]T, error) {
	var result []T
	data, err := f.send("GET", f.baseURL+"?"+query)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (f fetcher[_]) send(method string, uri string) ([]byte, error) {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
