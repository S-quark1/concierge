package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Service struct {
	ID          int64  `json:"id"`
	ServiceName string `json:"service_name"`
	ServiceType int64  `json:"service_type"`
	ServiceInfo string `json:"service_info"`
	Version     int32  `json:"version"`
}

type ServiceModel struct {
	DB *sql.DB
}

// admin can add new services which further would be presented in client's panel
func (s ServiceModel) Insert(service *Service) error {
	query := `
INSERT INTO service (service_name, service_type, service_info)
VALUES ($1, $2, $3)
RETURNING id, version`

	args := []interface{}{service.ServiceName, service.ServiceType, service.ServiceInfo}

	return s.DB.QueryRow(query, args...).Scan(&service.ID, &service.Version)
}

// get services
func (s ServiceModel) Get(id int64) (*Service, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
SELECT id, service_name, service_type, service_info, version
FROM service
WHERE id = $1`
	var service Service

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.ServiceName,
		&service.ServiceType,
		&service.ServiceInfo,
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
