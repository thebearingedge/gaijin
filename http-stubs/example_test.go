package example

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStubbedFetcher(t *testing.T) {
	want := "<h1>Example Domain</h1>"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, want)
	}))
	defer ts.Close()

	ex := Fetcher{url: ts.URL, client: *ts.Client()}
	got, err := ex.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	if got != want+"\n" {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestStubbedFetchFunc(t *testing.T) {
	want := "<h1>Example Domain</h1>"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, want)
	}))
	defer ts.Close()

	got, err := Fetch(*ts.Client(), ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if got != want+"\n" {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
