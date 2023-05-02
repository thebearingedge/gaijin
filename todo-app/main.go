package main

import (
	"database/sql"
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

	app.Group("/v1/todos").
		GET("", handler.GetAllTodos(repo)).
		GET("/:id", handler.GetOneTodoByID(repo)).
		POST("", handler.CreateOneTodo(repo)).
		PUT("/:id", handler.UpdateOneTodoByID(repo))

	return app
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	app := CreateApp(db)

	app.Use(
		gin.LoggerWithWriter(gin.DefaultWriter),
		gin.Recovery(),
	)

	gin.SetMode(gin.ReleaseMode)

	if err := app.Run(os.Getenv("LISTEN_ADDRESS")); err != nil {
		panic(err)
	}
}
