package repositories

import (
	"database/sql"
	"os"
	"testing"
	. "todo-app/data"

	"github.com/google/uuid"
)

func startTransaction(t *testing.T) *sql.Tx {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return tx
}

func TestEmptyTableGetsNoRows(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todos, err := r.GetAll()
	if err != nil {
		t.Fatalf("%v", err)
	}
	want := 0
	got := len(*todos)
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
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
	if err != nil {
		t.Fatalf("%v", err)
	}
	todos, err := r.GetAll()
	if err != nil {
		t.Fatalf("%v", err)
	}
	want := 3
	got := len(*todos)
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestEmptyTableGetsOneNilTodo(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todo, _ := r.GetOne(uuid.New())
	if todo != nil {
		t.Errorf("got %v, want %v", todo, nil)
	}
}

func TestEmptyTableUpdatesOneNilTodo(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todo, _ := r.UpdateOne(Todo{ID: uuid.New(), Task: ""})
	if todo != nil {
		t.Errorf("got %v, want %v", todo, nil)
	}
}

func TestCreatesOneTodo(t *testing.T) {
	r := NewTodosRepository(startTransaction(t))
	todo, err := r.CreateOne(Todo{Task: "Learn Go"})
	if err != nil {
		t.Fatalf("%v", err)
	}
	if todo.Task != "Learn Go" {
		t.Errorf("todo not created: got %v", todo)
	}
}

func TestNonEmptyTableGetsOneTodo(t *testing.T) {
	tx := startTransaction(t)
	r := NewTodosRepository(tx)
	id := uuid.New()
	_, err := tx.Exec(`--sql
		insert into "todos" ("todoId", "task")
		values ($1, 'Learn Go')
	`, id)
	if err != nil {
		t.Fatalf("%v", err)
	}
	todo, err := r.GetOne(id)
	if err != nil {
		t.Fatalf("todo not created: %v", err)
	}
	if todo.Task != "Learn Go" {
		t.Errorf("todo not created: got %v", todo)
	}
}

func TestNonEmptyTableUpdatesOneTodo(t *testing.T) {
	tx := startTransaction(t)
	r := NewTodosRepository(tx)
	id := uuid.New()
	_, err := tx.Exec(`--sql
		insert into "todos" ("todoId", "task")
		values ($1, 'Learn Go')
	`, id)
	if err != nil {
		t.Fatalf("%v", err)
	}
	todo, err := r.UpdateOne(Todo{ID: id, Task: "Accept Go"})
	if err != nil {
		t.Fatalf("todo not updated: %v", err)
	}
	if todo.Task != "Accept Go" {
		t.Errorf("todo not updated: got %v", todo)
	}
}
