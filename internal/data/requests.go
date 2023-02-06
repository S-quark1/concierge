package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Request struct {
	ID          int64          `json:"id"`
	ClientID    int64          `json:"client_id"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   sql.NullString `json:"deleted_at"`
	UpdatedAt   sql.NullString `json:"updated_at"`
}

type RequestModel struct {
	DB *sql.DB
}

func (r RequestModel) Insert(request *Request) error {
	query := `
INSERT INTO request (client_id, type, description, status, created_at) 
VALUES ($1, $2, $3, $4, NOW()) 
RETURNING id`
	args := []interface{}{request.ClientID, request.Type, request.Description, request.Status}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&request.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` ||
			err.Error() == `pq: повторяющееся значение ключа нарушает ограничение уникальности "users_email_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (r RequestModel) GetByID(id int64) (*Request, error) {
	query := `
SELECT * 
FROM request 
WHERE id = $1 AND deleted_at IS NULL`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var request Request

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.ClientID,
		&request.Type,
		&request.Description,
		&request.Status,
		&request.CreatedAt,
		&request.DeletedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &request, nil
}

func (m RequestModel) Update(request *Request) error {
	query := `
UPDATE request
SET type = $1, description = $2, status = $3, updated_at = $4
WHERE id = $5 AND deleted_at IS NULL`
	_, err := m.DB.Exec(query, request.Type, request.Description, request.Status, time.Now(), request.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (r RequestModel) Delete(id int64) error {
	query := `
DELETE FROM request
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r RequestModel) SoftDelete(id int64) error {
	query := `
UPDATE request
SET deleted_at = NOW() 
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
