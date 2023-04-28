package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	. "with-db/todos"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	t := Todo{Task: "Learn Go"}
	var created Todo
	err = db.QueryRow(`--sql
		insert into "todos" ("task")
		values ($1)
		returning "todoId",
		          "task",
				  "isCompleted",
				  "createdAt",
				  "updatedAt"
	`, t.Task).Scan(
		&created.TodoID,
		&created.Task,
		&created.IsCompleted,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error! %v\n", err)
		os.Exit(1)
	}
	data, err := json.Marshal(created)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error! %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("created: %+v\n", string(data))
}
