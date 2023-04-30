package handlers

import (
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

func (t TodosHandler) GetAllTodos(c *gin.Context)
func (t TodosHandler) GetOneTodoByID(c *gin.Context)
func (t TodosHandler) CreateOneTodo(c *gin.Context)
func (t TodosHandler) UpdateOneTodoByID(c *gin.Context)
