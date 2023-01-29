package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Company struct {
	ID        int64     `json:"id"`
	Code      int       `json:"code"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyModel struct {
	DB *sql.DB
}

func (c CompanyModel) Exists(id int) error {
	query := `SELECT id FROM company
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, id).Scan()
}

func (c *CompanyModel) Insert(company *Company) error {
	query := `
INSERT INTO company (code, name, full_name)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at`
	args := []interface{}{company.Code, company.Name, company.FullName}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
}

func (c CompanyModel) GetById(id int64) (*Company, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, code, name, full_name, created_at, deleted_at, updated_at
FROM company 
WHERE id = $1  AND deleted_at IS NULL`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var company Company

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.Code,
		&company.Name,
		&company.FullName,
		&company.CreatedAt,
		&company.DeletedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Otherwise, return a pointer to the Movie struct.
	return &company, nil
}

func (c *CompanyModel) Update(company *Company) error {
	query := `
UPDATE company
SET code = $1, name = $2, full_name = $3, updated_at = NOW()
WHERE id = $4 AND deleted_at IS NULL`
	args := []interface{}{company.Code, company.Name, company.FullName, company.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, args...)
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

func (c *CompanyModel) Delete(id int64) error {
	query := `DELETE FROM company WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := c.DB.ExecContext(ctx, query, id)
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

func (c *CompanyModel) SoftDelete(id int64) error {
	query := `UPDATE company SET deleted_at = NOW() WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := c.DB.ExecContext(ctx, query, id)
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
