package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
    const (
        insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
        confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
    )
	stmt, err := s.db.Prepare(insert)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(ctx, subject, description)
	if err != nil {
		log.Println("Error inserting TODO: ", err)
		return nil, err
	}
    id, err := result.LastInsertId()
    if err != nil {
        log.Println("Error getting last insert id: ", err)
        return nil, err
    }
    todo := &model.TODO{}
    res := s.db.QueryRowContext(ctx, confirm, id);
	todo.ID = id
    if err := res.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt ); err != nil {
        log.Println("Error scanning query result: ", err)
        return nil, err
    }
    return todo, nil
}
// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	var stmt *sql.Stmt
	var err error
	if prevID == 0{
		stmt, err = s.db.PrepareContext(ctx, read)
	} else {
		stmt, err = s.db.PrepareContext(ctx, readWithID)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var rows *sql.Rows
	if prevID == 0 {
		 rows, err = stmt.QueryContext(ctx, size)
		if err != nil {
			log.Println("Error reading TODOs: ", err)
			return nil, err
		}
		defer rows.Close()
	} else {
		rows, err = stmt.QueryContext(ctx, prevID, size)
		if err != nil {
			log.Println("Error reading TODOs: ", err)
			return nil, err
		}
		defer rows.Close()
	}
	todos := []*model.TODO{}
	for rows.Next() {
		todo := &model.TODO{}
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Println("Error scanning query result: ", err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows: ", err)
		return nil, err
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmt, err := s.db.PrepareContext(ctx,update)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		log.Println("Error updating TODO: ", err)
		return nil, err
	}
	res := s.db.QueryRowContext(ctx, confirm, id)

	todo := &model.TODO{}
	todo.ID = id
	if err := res.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt ); err != nil {
		return nil, &model.ErrNotFound{Message: err.Error()}
	}
	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	if len(ids) == 0{
		return &model.ErrNotFound{Message: "No ID to delete"}
	}
	placeholders := strings.Repeat(",?", len(ids)-1)
	deleteQuery := fmt.Sprintf(deleteFmt, placeholders)
	stmt, err := s.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return &model.ErrNotFound{Message: err.Error()}
	}
	defer stmt.Close()
	args := make([]interface{}, len(ids))
    for i, id := range ids {
        args[i] = id
    }
	res, err := stmt.ExecContext(ctx, args...)
    if err != nil {
        return &model.ErrNotFound{}
    }
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &model.ErrNotFound{}
	}
	
	return nil
}
