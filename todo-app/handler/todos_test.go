package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

func MustUnmarshal[T any](data []byte) T {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		panic(err)
	}
	return t
}

type StubGetAll struct {
	stub func() ([]Todo, error)
}

func (r StubGetAll) GetAll() ([]Todo, error) {
	return r.stub()
}

func TestGetAllTodosError(t *testing.T) {
	r := StubGetAll{func() ([]Todo, error) {
		return nil, errors.New("oops!")
	}}
	h := GetAllTodos(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAllTodosEmpty(t *testing.T) {
	want := make([]Todo, 0)
	r := StubGetAll{func() ([]Todo, error) {
		return want, nil
	}}
	h := GetAllTodos(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	got := MustUnmarshal[[]Todo](w.Body.Bytes())
	assert.Equal(t, want, got)
}

func TestGetAllTodosPopulated(t *testing.T) {
	want := []Todo{{Task: "Learn Go"}, {Task: "Accept Go"}}
	r := StubGetAll{func() ([]Todo, error) {
		return want, nil
	}}
	h := GetAllTodos(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	got := MustUnmarshal[[]Todo](w.Body.Bytes())
	assert.Equal(t, want, got)
}

type StubGetOneByID struct {
	stub func(id uuid.UUID) (*Todo, error)
}

func (r StubGetOneByID) GetOneByID(id uuid.UUID) (*Todo, error) {
	return r.stub(id)
}

func TestGetOneTodoNotFoundMissingIDParam(t *testing.T) {
	r := StubGetOneByID{func(id uuid.UUID) (*Todo, error) {
		return nil, nil
	}}
	h := GetOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOneTodoNotFoundInvalidIDParam(t *testing.T) {
	r := StubGetOneByID{func(id uuid.UUID) (*Todo, error) {
		return nil, nil
	}}
	h := GetOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asdf"})
	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOneTodoNotFoundValidIDParam(t *testing.T) {
	r := StubGetOneByID{func(id uuid.UUID) (*Todo, error) {
		return nil, nil
	}}
	h := GetOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOneTodoErrorValidIDParam(t *testing.T) {
	r := StubGetOneByID{func(id uuid.UUID) (*Todo, error) {
		return nil, errors.New("oops!")
	}}
	h := GetOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetOneTodoOkValidIDParam(t *testing.T) {
	want := Todo{ID: uuid.NewString()}
	r := StubGetOneByID{func(id uuid.UUID) (*Todo, error) {
		return &want, nil
	}}
	h := GetOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: want.ID})
	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	got := MustUnmarshal[Todo](w.Body.Bytes())
	assert.Equal(t, want, got)
}

type stubUpdateOneByID struct {
	stub func(id uuid.UUID, todo Todo) (*Todo, error)
}

func (r stubUpdateOneByID) UpdateOneByID(id uuid.UUID, todo Todo) (*Todo, error) {
	return r.stub(id, todo)
}

func TestUpdateOneTodoNotFoundMissingIDParam(t *testing.T) {
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOneTodoNotFoundInvalidIDParam(t *testing.T) {
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "asdf"})
	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOneTodoBadRequestNoBody(t *testing.T) {
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	h(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOneTodoBadRequestInvalidBody(t *testing.T) {
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString("{}"))
	h(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOneTodoNotFound(t *testing.T) {
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	h(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOneTodoError(t *testing.T) {
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return nil, errors.New("oops!")
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: uuid.NewString()})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateOneTodoOk(t *testing.T) {
	id := uuid.New()
	want := Todo{ID: id.String(), Task: "Accept Go"}
	r := stubUpdateOneByID{func(id uuid.UUID, todo Todo) (*Todo, error) {
		return &want, nil
	}}
	h := UpdateOneTodoByID(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	todo := Todo{ID: uuid.NewString(), Task: "Accept Go"}
	c.Params = append(c.Params, gin.Param{Key: "id", Value: todo.ID})
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Accept Go"}`))
	h(c)
	assert.Equal(t, http.StatusOK, w.Code)
	got := MustUnmarshal[Todo](w.Body.Bytes())
	assert.Equal(t, want, got)
}

type StubCreateOne struct {
	stub func(todo Todo) (*Todo, error)
}

func (r StubCreateOne) CreateOne(todo Todo) (*Todo, error) {
	return r.stub(todo)
}

func TestCreateOneTodoNoBody(t *testing.T) {
	r := StubCreateOne{func(todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := CreateOneTodo(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOneTodoInvalid(t *testing.T) {
	r := StubCreateOne{func(todo Todo) (*Todo, error) {
		return nil, nil
	}}
	h := CreateOneTodo(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString("{}"))
	h(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOneTodoError(t *testing.T) {
	r := StubCreateOne{func(todo Todo) (*Todo, error) {
		return nil, errors.New("oops!")
	}}
	h := CreateOneTodo(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"task":"Learn Go"}`))
	h(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateOneTodoOk(t *testing.T) {
	want := Todo{Task: "Accept Go"}
	r := StubCreateOne{func(todo Todo) (*Todo, error) {
		return &want, nil
	}}
	h := CreateOneTodo(r)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/", bytes.NewBufferString(`{"task":"Accept Go"}`))
	h(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	got := MustUnmarshal[Todo](w.Body.Bytes())
	assert.Equal(t, want, got)
}
