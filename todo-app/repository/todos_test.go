package repository

import (
	"database/sql"
	"os"
	"testing"
	. "todo-app/data"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func withRollback(t *testing.T, f func(*sql.Tx)) {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	assert.Nil(t, err)
	defer conn.Close()
	tx, err := conn.Begin()
	assert.Nil(t, err)
	defer tx.Rollback()
	f(tx)
}

func TestEmptyTableGetsNoRows(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		todos, err := r.GetAll()
		assert.Nil(t, err)
		assert.Len(t, todos, 0)
	})
}

func TestPopulatedTableGetsAllRows(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		_, err := tx.Exec(`--sql
			insert into "todos" ("task")
			values ('Learn Go'),
				   ('Do a Barrel Roll'),
				   ('Try a Somersault')
		`)
		assert.Nil(t, err)
		todos, err := r.GetAll()
		assert.Nil(t, err)
		assert.Len(t, todos, 3)
	})
}

func TestEmptyTableGetsOneNilTodo(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		todo, err := r.GetOneByID(uuid.New())
		assert.Nil(t, err)
		assert.Nil(t, todo)
	})
}

func TestNonEmptyTableGetsOneTodo(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		id := uuid.New()
		_, err := tx.Exec(`--sql
			insert into "todos" ("todoId", "task")
			values ($1, 'Learn Go')
		`, id)
		assert.Nil(t, err)
		todo, err := r.GetOneByID(id)
		assert.Nil(t, err)
		assert.Equal(t, todo.Task, "Learn Go")
	})
}

func TestEmptyTableUpdatesOneNilTodo(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		todo, err := r.UpdateOneByID(uuid.New(), Todo{Task: ""})
		assert.Nil(t, err)
		assert.Nil(t, todo)
	})
}

func TestNonEmptyTableUpdatesOneTodo(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		id := uuid.New()
		_, err := tx.Exec(`--sql
			insert into "todos" ("todoId", "task")
			values ($1, 'Learn Go')
		`, id)
		assert.Nil(t, err)
		todo, err := r.UpdateOneByID(id, Todo{Task: "Accept Go"})
		assert.Nil(t, err)
		assert.Equal(t, todo.Task, "Accept Go")
	})
}

func TestCreatesOneTodo(t *testing.T) {
	withRollback(t, func(tx *sql.Tx) {
		r := NewTodoRepository(tx)
		todo, err := r.CreateOne(Todo{Task: "Learn Go"})
		assert.Nil(t, err)
		assert.Equal(t, todo.Task, "Learn Go")
	})
}
