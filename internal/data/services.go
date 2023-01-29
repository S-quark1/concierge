package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Service struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedByID int64     `json:"created_by_id"`
	CompanyID   int       `json:"company_id"`
	CreatedAt   time.Time `json:"created_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ServiceModel struct {
	DB *sql.DB
}

// cs employee can add new services which further would be presented in client's panel
func (s *ServiceModel) Insert(service *Service) error {
	query := `
INSERT INTO service (name, description, type, created_by_id, company_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at`
	args := []interface{}{service.Name, service.Description, service.Type, service.CreatedByID, service.CompanyID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID, &service.CreatedAt, &service.UpdatedAt)
}

// get services, companies, employee and prices
func (s *ServiceModel) GetById(id int64) (*Service, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, name, description, type, created_by_id, company_id, created_at, deleted_at, updated_at
FROM service
WHERE id = $1  AND deleted_at IS NULL`
	var service Service
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.Type,
		&service.CreatedByID,
		&service.CompanyID,
		&service.CreatedAt,
		&service.DeletedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &service, nil
}

func (s *ServiceModel) Update(service *Service) error {
	query := `
UPDATE service
SET name = $1, description = $2, type = $3, created_by_id = $4, company_id = $5, updated_at = NOW()
WHERE id = $6 AND deleted_at IS NULL`
	args := []interface{}{service.Name, service.Description, service.Type, service.CreatedByID, service.CompanyID, service.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := s.DB.ExecContext(ctx, query, args...)
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

func (s *ServiceModel) Delete(id int64) error {
	query := `DELETE FROM service WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := s.DB.ExecContext(ctx, query, id)
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

func (s *ServiceModel) SoftDelete(id int64) error {
	query := `
UPDATE service
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, id)
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

type Price struct {
	ID        int64     `json:"id"`
	ServiceID int64     `json:"service_id"`
	Price     int       `json:"price"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PriceModel struct {
	DB *sql.DB
}

func (p *PriceModel) Insert(price *Price) error { // TODO: проверить существует ли id
	query := `
INSERT INTO price (service_id, price, user_type)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at`
	args := []interface{}{price.ServiceID, price.Price, price.UserType}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&price.ID, &price.CreatedAt, &price.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// get all prices for a specific service
func (p *PriceModel) GetByServiceId(serviceID int64) ([]*Price, error) {
	if serviceID < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, service_id, price, user_type, created_at, deleted_at, updated_at
FROM price
WHERE service_id = $1 AND deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	prices := []*Price{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var price Price
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&price.ID,
			&price.ServiceID,
			&price.Price,
			&price.UserType,
			&price.CreatedAt,
			&price.DeletedAt,
			&price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Add the Movie struct to the slice.
		prices = append(prices, &price)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return prices, nil
}

func (p *PriceModel) Update(price *Price) error {
	query := `
UPDATE price
SET price = $1, user_type = $2, service_id = $3, updated_at = NOW()
WHERE id = $4 AND deleted_at IS NULL`
	args := []interface{}{price.Price, price.UserType, price.ServiceID, price.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&price.UpdatedAt)
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

func (p *PriceModel) Delete(id int64) error {
	query := `DELETE FROM price WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := p.DB.ExecContext(ctx, query, id)
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

func (p *PriceModel) SoftDelete(id int64) error {
	query := `
UPDATE price
SET deleted_at = NOW()
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := p.DB.ExecContext(ctx, query, id)
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
