package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Service struct {
	/*id bigserial PRIMARY KEY,
	name text NOT NULL,
	description text NOT NULL,
	dateOfCreation timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	type text NOT NULL,
	// type_id int NOT NULL,   // actually no, I don't need them, just make a table with entry service_id (k) and type text
	createdBy_id bigint NOT NULL,	//k
	company_id int NOT NULL,
	//prices_id bigint NOT NULL,  // actually no, I don't need them, just make a table with entry service_id (k) and type text
	version int NOT NULL DEFAULT 1*/

	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	DateOfCreation    time.Time `json:"-"`
	ServiceType       string    `json:"type"`
	CreatedByEmployee int64     `json:"createdBy_id"`
	CompanyProviding  int32     `json:"company_id"`
	Version           int32     `json:"version"`
}

// Todo I need companies, cs employees and prices
type ServiceModel struct {
	DB *sql.DB
}

// cs employee can add new services which further would be presented in client's panel
func (s ServiceModel) Insert(service *Service) error {
	query := `
INSERT INTO service (name, description, type, createdBy_id, company_id)
VALUES ($1, $2, $3, $4, $5)
-- WHERE EXISTS (SELECT 1 FROM company WHERE id = $5)
RETURNING id, dateOfCreation, version`

	args := []interface{}{service.Name, service.Description, service.ServiceType, service.CreatedByEmployee, service.CompanyProviding}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID, &service.DateOfCreation, &service.Version)
}

// get services, companies, employee and prices
func (s ServiceModel) GetById(id int64) (*Service, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, name, description, dateOfCreation, type, createdBy_id, company_id, version
FROM service
WHERE id = $1`
	var service Service

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.DateOfCreation,
		&service.ServiceType,
		&service.CreatedByEmployee,
		&service.CompanyProviding,
		&service.Version,
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

type Price struct {
	/*service_id bigint NOT NULL,
	price int NOT NULL,
	user_type text NOT NULL,   // later create k
	version int NOT NULL DEFAULT 1*/
	ServiceID int64  `json:"service_id,omitempty"`
	Price     int32  `json:"price"`
	UserType  string `json:"user_type"`
	Version   int32  `json:"version"`
}

type PriceModel struct {
	DB *sql.DB
}

func (p PriceModel) Insert(price *Price) error { // TODO: проверить существует ли id
	query := `
INSERT INTO price (service_id, price, user_type)
VALUES ($1, $2, $3)
RETURNING version`

	args := []interface{}{price.ServiceID, price.Price, price.UserType}

	return p.DB.QueryRow(query, args...).Scan(&price.Version)
}

// get all prices for a specific service
func (p PriceModel) GetById(serviceID int64) ([]*Price, error) {
	if serviceID < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT price, user_type, version
FROM price
WHERE service_id = $1`

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
			&price.Price,
			&price.UserType,
			&price.Version,
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
