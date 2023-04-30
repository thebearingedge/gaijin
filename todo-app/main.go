package main

import (
	"database/sql"
	"fmt"
	"os"
	. "todo-app/handlers"
	. "todo-app/repositories"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to database - %v", err)
		os.Exit(1)
	}

	server := gin.Default()

	todosHandler := NewTodosHandler(NewTodosRepository(db))

	todosRouter := server.Group("/v1/todos")
	todosRouter.GET("/", todosHandler.GetAllTodos)
	todosRouter.GET("/:id", todosHandler.GetOneTodoByID)
	todosRouter.POST("/", todosHandler.CreateOneTodo)
	todosRouter.PUT("/:id", todosHandler.UpdateOneTodoByID)

	if err := server.Run(os.Getenv("LISTEN_ADDRESS")); err != nil {
		fmt.Fprintf(os.Stderr, "could start the server - %v", err)
		os.Exit(1)
	}
}
