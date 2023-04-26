package example

import (
	"io"
	"net/http"
)

type Fetcher struct {
	url    string
	client http.Client
}

func (f Fetcher) Fetch() (string, error) {
	resp, err := f.client.Get(f.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), nil
}

// just a function
func Fetch(client http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), nil
}
