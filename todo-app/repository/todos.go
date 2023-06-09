package repository

import (
	"database/sql"
	"fmt"
	. "todo-app/data"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TodoRepository struct {
	db DB
}

func NewTodoRepository(db DB) *TodoRepository {
	return &TodoRepository{db}
}

func (r *TodoRepository) CreateOne(t Todo) (*Todo, error) {
	var todo Todo
	row := r.db.QueryRow(`--sql
		insert into "todos" ("task", "isCompleted")
		values ($1, $2)
		returning "todoId",
				  "task",
				  "isCompleted",
				  "createdAt",
				  "updatedAt"
	`, t.Task, t.IsCompleted)
	err := row.Scan(&todo.ID, &todo.Task, &todo.IsCompleted, &todo.CreatedAt, &todo.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}
	return &todo, nil
}

func (r *TodoRepository) UpdateOneByID(id uuid.UUID, t Todo) (*Todo, error) {
	var todo Todo
	row := r.db.QueryRow(`--sql
		update "todos"
		   set "task"        = coalesce($1, "task"),
			   "isCompleted" = coalesce($2, "isCompleted"),
			   "updatedAt"   = now()
		 where "todoId"      = $3
		returning "todoId",
				  "task",
				  "isCompleted",
				  "createdAt",
				  "updatedAt"
	`, t.Task, t.IsCompleted, id)
	err := row.Scan(&todo.ID, &todo.Task, &todo.IsCompleted, &todo.CreatedAt, &todo.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}
	return &todo, nil
}

func (r *TodoRepository) GetOneByID(id uuid.UUID) (*Todo, error) {
	var todo Todo
	row := r.db.QueryRow(`--sql
		select "todoId",
			   "task",
			   "isCompleted",
			   "createdAt",
			   "updatedAt"
		  from "todos"
		 where "todoId" = $1
	`, id)
	err := row.Scan(&todo.ID, &todo.Task, &todo.IsCompleted, &todo.CreatedAt, &todo.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}
	return &todo, nil
}

func (r *TodoRepository) GetAll() ([]Todo, error) {
	rows, err := r.db.Query(`--sql
		select "todoId",
			   "task",
			   "isCompleted",
			   "createdAt",
			   "updatedAt"
		  from "todos"
	`)
	if err != nil {
		return nil, fmt.Errorf("querying database: %w", err)
	}
	all := make([]Todo, 0)
	defer rows.Close()
	for rows.Next() {
		todo := Todo{}
		err := rows.Scan(
			&todo.ID, &todo.Task, &todo.IsCompleted, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		all = append(all, todo)
	}
	return all, nil
}
