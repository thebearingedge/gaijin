package repository

import (
	"database/sql"
	"os"
	"testing"
	. "todo-app/data"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func startTransaction(t *testing.T) *sql.Tx {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	assert.Nil(t, err)
	tx, err := conn.Begin()
	assert.Nil(t, err)
	return tx
}

func TestEmptyTableGetsNoRows(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todos, err := r.GetAll()
	assert.Nil(t, err)
	assert.Len(t, *todos, 0)
}

func TestPopulatedTableGetsAllRows(t *testing.T) {
	tx := startTransaction(t)
	r := NewTodosRepository(tx)
	_, err := tx.Exec(`--sql
		insert into "todos" ("task")
		values ('Learn Go'),
			   ('Do a Barrel Roll'),
			   ('Try a Somersault')
	`)
	assert.Nil(t, err)
	todos, err := r.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, len(*todos), 3)
}

func TestEmptyTableGetsOneNilTodo(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todo, err := r.GetOne(uuid.New())
	assert.Nil(t, err)
	assert.Nil(t, todo)
}

func TestEmptyTableUpdatesOneNilTodo(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todo, err := r.UpdateOne(uuid.New(), Todo{Task: ""})
	assert.Nil(t, err)
	assert.Nil(t, todo)
}

func TestCreatesOneTodo(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todo, err := r.CreateOne(Todo{Task: "Learn Go"})
	assert.Nil(t, err)
	assert.Equal(t, todo.Task, "Learn Go")
}

func TestNonEmptyTableGetsOneTodo(t *testing.T) {
	tx := startTransaction(t)
	r := NewTodosRepository(tx)
	id := uuid.New()
	_, err := tx.Exec(`--sql
		insert into "todos" ("todoId", "task")
		values ($1, 'Learn Go')
	`, id)
	assert.Nil(t, err)
	todo, err := r.GetOne(id)
	assert.Nil(t, err)
	assert.Equal(t, todo.Task, "Learn Go")

}

func TestNonEmptyTableUpdatesOneTodo(t *testing.T) {
	tx := startTransaction(t)
	r := NewTodosRepository(tx)
	id := uuid.New()
	_, err := tx.Exec(`--sql
		insert into "todos" ("todoId", "task")
		values ($1, 'Learn Go')
	`, id)
	assert.Nil(t, err)
	todo, err := r.UpdateOne(id, Todo{Task: "Accept Go"})
	assert.Nil(t, err)
	assert.Equal(t, todo.Task, "Accept Go")
}
