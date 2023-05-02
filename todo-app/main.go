package main

import (
	"database/sql"
	"fmt"
	"os"
	"todo-app/data"
	"todo-app/handler"
	. "todo-app/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func CreateApp(db data.DB) *gin.Engine {
	app := gin.New()

	repo := NewTodosRepository(db)

	todoHandler := handler.NewTodosHandler(repo)

	todoRoutes := app.Group("/v1/todos")

	todoRoutes.GET("", handler.GetAllTodos(repo))
	todoRoutes.GET("/:id", handler.GetOneTodoById(repo))
	todoRoutes.POST("", todoHandler.CreateOneTodo)
	todoRoutes.PUT("/:id", handler.UpdateOneByID(repo))

	return app
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	app := CreateApp(db)

	app.Use(
		gin.LoggerWithWriter(gin.DefaultWriter),
		gin.Recovery(),
	)

	if err := app.Run(os.Getenv("LISTEN_ADDRESS")); err != nil {
		fmt.Fprintf(os.Stderr, "could start the app - %v", err)
		os.Exit(1)
	}
}
