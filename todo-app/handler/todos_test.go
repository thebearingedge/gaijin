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

type StubGetAll struct {
	stub func() ([]Todo, error)
}

func (r StubGetAll) GetAll() ([]Todo, error) {
	return r.stub()
}

func stubGetAll(f func() ([]Todo, error)) getAll {
	return StubGetAll{f}
}

func TestGetAllTodosError(t *testing.T) {
	r := stubGetAll(func() ([]Todo, error) {
		return nil, errors.New("oops!")
	})
	h := GetAllTodos(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.Equal(t, w.Body.String(), "")
}

func TestGetAllTodosEmpty(t *testing.T) {
	r := stubGetAll(func() ([]Todo, error) {
		return make([]Todo, 0), nil
	})
	h := GetAllTodos(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, w.Code, http.StatusOK)
	want, _ := r.GetAll()
	var got []Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestGetAllTodosPopulated(t *testing.T) {
	r := stubGetAll(func() ([]Todo, error) {
		return []Todo{{Task: "Learn Go"}, {Task: "Accept Go"}}, nil
	})
	h := GetAllTodos(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	want, _ := r.GetAll()
	var got []Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

type StubGetOneByID struct {
	stub func(id uuid.UUID) (*Todo, error)
}

func (r StubGetOneByID) GetOneByID(id uuid.UUID) (*Todo, error) {
	return r.stub(id)
}

func stubGetOneByID(f func(uuid.UUID) (*Todo, error)) getOneByID {
	return StubGetOneByID{f}
}

func TestGetOneTodoNotFoundMissingIDParam(t *testing.T) {
	handler := GetOneTodoByID(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoNotFoundInvalidIDParam(t *testing.T) {
	handler := GetOneTodoByID(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asdf"})
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoNotFoundValidIDParam(t *testing.T) {
	handler := GetOneTodoByID(&GetOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestGetOneTodoErrorValidIDParam(t *testing.T) {
	handler := GetOneTodoByID(&GetOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestGetOneTodoOkValidIDParam(t *testing.T) {
	handler := GetOneTodoByID(&GetOneTodoOk{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id := uuid.New()
	c.Params = append(c.Params, gin.Param{Key: "id", Value: id.String()})
	handler(c)
	assert.Equal(t, w.Code, http.StatusOK)
	want, _ := (&GetOneTodoOk{}).GetOneByID(id)
	var got Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, &got)
}

type UpdateOneTodoNotFound struct{}

func (t UpdateOneTodoNotFound) UpdateOneByID(id uuid.UUID, todo Todo) (*Todo, error) {
	return nil, nil
}

func TestUpdateOneTodoNotFoundMissingIDParam(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, w.Code, http.StatusNotFound)
}

func TestUpdateOneTodoNotFoundInvalidIDParam(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asdf"})
	handler(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOneTodoBadRequestNoBody(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	handler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOneTodoBadRequestInvalidBody(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString("{}"))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOneTodoNotFound(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoNotFound{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

type UpdateOneTodoError struct{}

func (r *UpdateOneTodoError) UpdateOneByID(id uuid.UUID, todo Todo) (*Todo, error) {
	return nil, errors.New("oops!")
}

func TestUpdateOneTodoError(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type UpdateOneTodoOk struct{}

func (r *UpdateOneTodoOk) UpdateOneByID(id uuid.UUID, todo Todo) (*Todo, error) {
	return &Todo{ID: id.String(), Task: "Accept Go"}, nil
}

func TestUpdateOneTodoOk(t *testing.T) {
	handler := UpdateOneTodoByID(&UpdateOneTodoOk{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	todo := Todo{ID: uuid.NewString(), Task: "Accept Go"}
	c.Params = append(c.Params, gin.Param{Key: "id", Value: todo.ID})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Accept Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusOK, w.Code)
	want, _ := (&UpdateOneTodoOk{}).UpdateOneByID(uuid.MustParse(todo.ID), todo)
	var got Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, &got)
}

type CreateOneTodoError struct{}

func (r *CreateOneTodoError) CreateOne(todo Todo) (*Todo, error) {
	return nil, errors.New("oops!")
}

func TestCreateOneTodoNoBody(t *testing.T) {
	handler := CreateOneTodo(&CreateOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOneTodoInvalid(t *testing.T) {
	handler := CreateOneTodo(&CreateOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString("{}"))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOneTodoError(t *testing.T) {
	handler := CreateOneTodo(&CreateOneTodoError{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), "")
}

type CreateOneTodoOk struct{}

func (r *CreateOneTodoOk) CreateOne(todo Todo) (*Todo, error) {
	return &Todo{Task: "Accept Go"}, nil
}

func TestCreateOneTodoOk(t *testing.T) {
	handler := CreateOneTodo(&CreateOneTodoOk{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Accept Go"}`))
	c.Request.Header.Add("content-type", "application/json")
	handler(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	want := Todo{Task: "Accept Go"}
	var got Todo
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}
