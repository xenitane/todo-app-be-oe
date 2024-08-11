package database

import (
	"database/sql"
	"errors"

	"github.com/xenitane/todo-app-be-oe/internals/todo"
)

func (s *service) GetAllTodosForUser(userID int64) ([]*todo.Todo, error) {
	query := `select * from todos where owner_id = $1`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	todos := []*todo.Todo{}
	for rows.Next() {
		todo, err := scanTodoRow(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *service) GetTodoByIDForUser(tid, uid int64) (*todo.Todo, error) {
	query := `select * from todos where id = $1 and owner_id = $2;`
	rows, err := s.db.Query(query, tid, uid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanTodoRow(rows)
	}
	return nil, sql.ErrNoRows
}

func scanTodoRow(rows *sql.Rows) (*todo.Todo, error) {
	todo := new(todo.Todo)
	err := rows.Scan(
		&todo.TodoId,
		&todo.OwnerId,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.DueDate,
		&todo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *service) InsertTodo(t *todo.Todo) error {
	insertQuery := `insert into todos (owner_id, title, description, status, due_date) values ($1, $2, $3, $4, $5) returning id;`
	rows, err := s.db.Query(
		insertQuery,
		t.OwnerId,
		t.Title,
		t.Description,
		t.Status,
		t.DueDate,
	)
	if err != nil {
		return err
	}
	for rows.Next() {
		return rows.Scan(&t.TodoId)
	}
	return errors.New("could not insert")
}

func (s *service) DeleteTodoByIDForUser(tid, uid int64) error {
	deleteQuery := `delete from todos wherer id = $1 and owner_id = $2`
	res, err := s.db.Exec(deleteQuery, tid, uid)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil || ra == 0 {
		return err
	}
	return nil
}

func (s *service) UpdateTodoByIdForUser(t *todo.Todo) error {
	updateQry := `update todos set (title, description, status, due_date) = ($3, $4, $5, $6) where id = $1 and owner_id = $2`
	_, err := s.db.Exec(updateQry, t.TodoId, t.OwnerId, t.Title, t.Description, t.Status, t.DueDate)
	return err
}
