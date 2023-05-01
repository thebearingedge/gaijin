package handlers

import (
	"net/http"
	. "todo-app/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Todos interface {
	GetAll() (*[]Todo, error)
	GetOne(id uuid.UUID) (*Todo, error)
	UpdateOne(t Todo) (*Todo, error)
	CreateOne(t Todo) (*Todo, error)
}

type TodosHandler struct {
	todos Todos
}

func NewTodosHandler(t Todos) TodosHandler {
	return TodosHandler{t}
}

func (t TodosHandler) GetAllTodos(c *gin.Context) {
	todos, err := t.todos.GetAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (t TodosHandler) GetOneTodoByID(c *gin.Context) {
	rawId, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	id, err := uuid.FromBytes([]byte(rawId))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	todo, err := t.todos.GetOne(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if todo == nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.JSON(http.StatusOK, todo)
}
func (t TodosHandler) CreateOneTodo(c *gin.Context)     {}
func (t TodosHandler) UpdateOneTodoByID(c *gin.Context) {}
