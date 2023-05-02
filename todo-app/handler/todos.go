package handler

import (
	"net/http"
	"strings"
	. "todo-app/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	GetOneByID(id uuid.UUID) (*Todo, error)
}

func GetOneTodoByID(t getOneById) func(c *gin.Context) {
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
		todo, err := t.GetOneByID(todoId)
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

type createOne interface {
	CreateOne(todo Todo) (*Todo, error)
}

func CreateOneTodo(t createOne) func(c *gin.Context) {
	return func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if &todo.Task == nil || strings.TrimSpace(todo.Task) == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		created, err := t.CreateOne(todo)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, created)
	}
}

type updateOne interface {
	UpdateOne(id uuid.UUID, todo Todo) (*Todo, error)
}

func UpdateOneTodoByID(t updateOne) func(c *gin.Context) {
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
