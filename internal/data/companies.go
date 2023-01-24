package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Company struct {
	/*id bigint PRIMARY KEY
	code int NOT NULL,
	name text NOT NULL,
	full_name text NOT NULL,
	version int NOT NULL DEFAULT 1*/

	ID       int64  `json:"id"`
	Code     int32  `json:"code"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Version  int32  `json:"version"`
}

type CompanyModel struct {
	DB *sql.DB
}

func (c CompanyModel) Exists(id int32) error {
	query := `SELECT id FROM company
WHERE id = $1
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, id).Scan()
}
func (c CompanyModel) GetById(id int64) (*Company, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, code, name, full_name, version
FROM company
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var company Company

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.Code,
		&company.Name,
		&company.FullName,
		&company.Version,
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
