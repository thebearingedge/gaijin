package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	. "todo-app/data"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type GetAllTodosError struct{}

func (r *GetAllTodosError) GetAll() (*[]Todo, error) {
	return nil, errors.New("oops!")
}

func TestGetAllTodosError(t *testing.T) {
	handler := GetAllTodos(&GetAllTodosError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.Equal(t, w.Body.String(), "")
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
	want, _ := (&GetAllTodosEmpty{}).GetAll()
	var got []Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, &got)
}

type GetAllTodosPopulated struct{}

func (r *GetAllTodosPopulated) GetAll() (*[]Todo, error) {
	todos := []Todo{
		{Task: "Learn Go"},
		{Task: "Accept Go"},
	}
	return &todos, nil
}

func TestGetAllTodosPopulated(t *testing.T) {
	handler := GetAllTodos(&GetAllTodosPopulated{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, http.StatusOK, w.Code)
	want, _ := (&GetAllTodosPopulated{}).GetAll()
	var got []Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, &got)
}

type GetOneTodoNotFound struct{}

func (r *GetOneTodoNotFound) GetOne(id uuid.UUID) (*Todo, error) {
	return nil, nil
}

func TestGetOneTodoNotFoundMissingIDParam(t *testing.T) {
	handler := GetOneTodoById(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoNotFoundInvalidIDParam(t *testing.T) {
	handler := GetOneTodoById(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asdf"})
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoNotFoundValidIDParam(t *testing.T) {
	handler := GetOneTodoById(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

type GetOneTodoError struct{}

func (r *GetOneTodoError) GetOne(id uuid.UUID) (*Todo, error) {
	return nil, errors.New("oops!")
}

func TestGetOneTodoErrorValidIDParam(t *testing.T) {
	handler := GetOneTodoById(&GetOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

type GetOneTodoOk struct{}

func (r *GetOneTodoOk) GetOne(id uuid.UUID) (*Todo, error) {
	return &Todo{Task: "Learn Go"}, nil
}

func TestGetOneTodoOkValidIDParam(t *testing.T) {
	handler := GetOneTodoById(&GetOneTodoOk{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id := uuid.New()
	c.Params = append(c.Params, gin.Param{Key: "id", Value: id.String()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusOK)
	want, _ := (&GetOneTodoOk{}).GetOne(id)
	var got Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, &got)
}

type UpdateOneTodoNotFound struct{}

func (t UpdateOneTodoNotFound) UpdateOne(id uuid.UUID, todo Todo) (*Todo, error) {
	return nil, nil
}

func TestUpdateOneTodoNotFoundMissingIDParam(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestUpdateOneTodoNotFoundInvalidIDParam(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asdf"})
	handler(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOneTodoBadRequestNoBody(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOneTodoBadRequestInvalidBody(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString("{}"))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOneTodoNotFound(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

type UpdateOneTodoError struct{}

func (r *UpdateOneTodoError) UpdateOne(id uuid.UUID, todo Todo) (*Todo, error) {
	return nil, errors.New("oops!")
}

func TestUpdateOneTodoError(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type UpdateOneTodoOk struct{}

func (r *UpdateOneTodoOk) UpdateOne(id uuid.UUID, todo Todo) (*Todo, error) {
	return &Todo{ID: id.String(), Task: "Accept Go"}, nil
}

func TestUpdateOneTodoOk(t *testing.T) {
	handler := UpdateOneByID(&UpdateOneTodoOk{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	todo := Todo{ID: uuid.NewString(), Task: "Accept Go"}
	c.Params = append(c.Params, gin.Param{Key: "id", Value: todo.ID})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Accept Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusOK, w.Code)
	want, _ := (&UpdateOneTodoOk{}).UpdateOne(uuid.MustParse(todo.ID), todo)
	var got Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, &got)
}
