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

func CreateApp(db *sql.DB) *gin.Engine {
	app := gin.Default()

	todosHandler := NewTodosHandler(NewTodosRepository(db))

	todosRouter := app.Group("/v1/todos")

	todosRouter.GET("/", todosHandler.GetAllTodos)
	todosRouter.GET("/:id", todosHandler.GetOneTodoByID)
	todosRouter.POST("/", todosHandler.CreateOneTodo)
	todosRouter.PUT("/:id", todosHandler.UpdateOneTodoByID)

	return app
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to database - %v", err)
		os.Exit(1)
	}

	app := CreateApp(db)

	if err := app.Run(os.Getenv("LISTEN_ADDRESS")); err != nil {
		fmt.Fprintf(os.Stderr, "could start the app - %v", err)
		os.Exit(1)
	}
}
