package repository

import (
	"database/sql"
	"fmt"
	. "todo-app/data"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TodosRepository struct {
	db DB
}

func NewTodosRepository(db DB) *TodosRepository {
	return &TodosRepository{db}
}

func (r *TodosRepository) CreateOne(t Todo) (*Todo, error) {
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
		return nil, fmt.Errorf("could not scan row - %w", err)
	}
	return &todo, nil
}

func (r *TodosRepository) UpdateOne(id uuid.UUID, t Todo) (*Todo, error) {
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
		return nil, fmt.Errorf("could not scan row - %w", err)
	}
	return &todo, nil
}

func (r *TodosRepository) GetOne(id uuid.UUID) (*Todo, error) {
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
		return nil, fmt.Errorf("could not scan row - %w", err)
	}
	return &todo, nil
}

func (r *TodosRepository) GetAll() (*[]Todo, error) {
	rows, err := r.db.Query(`--sql
		select "todoId",
			   "task",
			   "isCompleted",
			   "createdAt",
			   "updatedAt"
		  from "todos"
	`)
	if err != nil {
		return nil, fmt.Errorf("could not query database - %w", err)
	}
	todos := make([]Todo, 0)
	defer rows.Close()
	for rows.Next() {
		todo := Todo{}
		err := rows.Scan(
			&todo.ID, &todo.Task, &todo.IsCompleted, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan rows - %w", err)
		}
		todos = append(todos, todo)
	}
	return &todos, nil
}
