package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	savetodo, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	todoid, err := savetodo.LastInsertId()
	if err != nil {
		return nil, err
	}

	var Subject string
	var Description string
	var CreatedAt time.Time
	var UpdatedAt time.Time

	result := s.db.QueryRowContext(ctx, confirm, todoid)
	if err = result.Scan(&Subject, &Description, &CreatedAt, &UpdatedAt); err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          todoid,
		Subject:     Subject,
		Description: Description,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
	}

	return todo, err
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	todos := []*model.TODO{}
	var readSize int64
	if size == 0 {
		readSize = 3
	} else {
		readSize = size
	}

	if prevID != 0 {
		rows, err := s.db.QueryContext(ctx, readWithID, prevID, readSize)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var Id int64
			var Subject string
			var Description string
			var CreatedAt time.Time
			var UpdatedAt time.Time
			if err = rows.Scan(&Id, &Subject, &Description, &CreatedAt, &UpdatedAt); err != nil {
				return nil, err
			}
			todo := &model.TODO{
				ID:          Id,
				Subject:     Subject,
				Description: Description,
				CreatedAt:   CreatedAt,
				UpdatedAt:   UpdatedAt,
			}

			todos = append(todos, todo)
		}
		return todos, err
	}

	rows, err := s.db.QueryContext(ctx, read, readSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var Id int64
		var Subject string
		var Description string
		var CreatedAt time.Time
		var UpdatedAt time.Time
		if err = rows.Scan(&Id, &Subject, &Description, &CreatedAt, &UpdatedAt); err != nil {
			return nil, err
		}
		todo := &model.TODO{
			ID:          Id,
			Subject:     Subject,
			Description: Description,
			CreatedAt:   CreatedAt,
			UpdatedAt:   UpdatedAt,
		}
		todos = append(todos, todo)
	}
	return todos, err
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	_, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	var Subject string
	var Description string
	var CreatedAt time.Time
	var UpdatedAt time.Time

	result := s.db.QueryRowContext(ctx, confirm, id)
	if err = result.Scan(&Subject, &Description, &CreatedAt, &UpdatedAt); err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     Subject,
		Description: Description,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt,
	}

	return todo, err
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	if len(ids) == 0 {
		return nil
	}

	var arg []interface{}
	for _, id := range ids {
		idToInt := strconv.FormatInt(id, 10)
		arg = append(arg, idToInt)
	}

	var query string
	if len(ids) == 1 {
		query = fmt.Sprintf(deleteFmt, "")
	} else {
		querySymbol := strings.Repeat(string(',')+string('?'), len(ids)-1)
		query = fmt.Sprintf(deleteFmt, querySymbol)
	}

	result, err := s.db.ExecContext(ctx, query, arg...)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return model.ErrNotFound{}
	}

	return err
}
