package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorldRoute(t *testing.T) {
	app := createApp()

	r := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/hello", nil)
	assert.Nil(t, err)

	app.ServeHTTP(r, req)

	assert.Equal(t, 200, r.Code)
	assert.Equal(t, "world", r.Body.String())
}
