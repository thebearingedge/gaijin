package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func buildApp(t *testing.T) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	assert.Nil(t, err)
	tx, err := conn.Begin()
	assert.Nil(t, err)
	app := CreateApp(tx)
	return app
}

func TestGetAllTodosEmpty(t *testing.T) {
	app := buildApp(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/todos", nil)
	app.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.JSONEq(t, w.Body.String(), `[]`)
}

func TestGetOneTodoEmpty(t *testing.T) {
	app := buildApp(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/todos/"+uuid.New().String(), nil)
	app.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)
}
