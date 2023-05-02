package handler

import (
	"net/http"
	"strings"
	. "todo-app/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Todos interface {
	GetAll() (*[]Todo, error)
	GetOne(id uuid.UUID) (*Todo, error)
	UpdateOne(id uuid.UUID, t Todo) (*Todo, error)
	CreateOne(t Todo) (*Todo, error)
}

type TodosHandler struct {
	todos Todos
}

func NewTodosHandler(t Todos) TodosHandler {
	return TodosHandler{t}
}

type getAll interface {
	GetAll() (*[]Todo, error)
}

func GetAllTodos(t getAll) func(c *gin.Context) {
	return func(c *gin.Context) {
		todos, err := t.GetAll()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

type getOneById interface {
	GetOne(id uuid.UUID) (*Todo, error)
}

func GetOneTodoById(t getOneById) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		todoId, err := uuid.Parse(id)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		todo, err := t.GetOne(todoId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if todo == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

func (t TodosHandler) CreateOneTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if &todo.Task == nil || strings.TrimSpace(todo.Task) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	created, err := t.todos.CreateOne(todo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, created)
}

type UpdateOne interface {
	UpdateOne(id uuid.UUID, todo Todo) (*Todo, error)
}

func UpdateOneByID(t UpdateOne) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		todoId, err := uuid.Parse(id)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if &todo.Task == nil || strings.TrimSpace(todo.Task) == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		updated, err := t.UpdateOne(todoId, todo)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if updated == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, updated)
	}
}
