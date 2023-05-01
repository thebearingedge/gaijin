package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	. "todo-app/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type GetAllTodosEmpty struct{}

func (r *GetAllTodosEmpty) GetAll() (*[]Todo, error) {
	todos := make([]Todo, 0)
	return &todos, nil
}

func TestGetAllTodosEmpty(t *testing.T) {
	handler := GetAllTodos(&GetAllTodosEmpty{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.String(), `[]`)
}

type GetOneTodoNotFound struct{}

func (r *GetOneTodoNotFound) GetOne(id uuid.UUID) (*Todo, error) {
	return nil, nil
}

func TestGetOneTodoNotFoundMissingIDParam(t *testing.T) {
	handler := GetOneTodo(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoNotFoundInvalidIDParam(t *testing.T) {
	handler := GetOneTodo(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "yo-mama"})
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoNotFoundValidIDParam(t *testing.T) {
	handler := GetOneTodo(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}
