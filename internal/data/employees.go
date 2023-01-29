package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/concierge/service/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)

/*id bigint PRIMARY KEY,
firstName text NOT NULL,
lastName text NOT NULL,
username text UNIQUE NOT NULL,
password_hash bytea NOT NULL,
role smallint NOT NULL,
dateOfCreation timestamp(0) with time zone NOT NULL DEFAULT NOW(),*/

type CSEmployee struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  password  `json:"password_hash"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CSEmployeeModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(validator.Matches(username, validator.UsernameRX), "username", "must be a valid username")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateEmployee(v *validator.Validator, employee *CSEmployee) {
	v.Check(employee.FirstName != "", "firstname", "must be provided")
	v.Check(len(employee.FirstName) <= 500, "firstname", "must not be more than 500 bytes long")

	v.Check(employee.LastName != "", "lastname", "must be provided")
	v.Check(len(employee.LastName) <= 500, "lastname", "must not be more than 500 bytes long")

	ValidateUsername(v, employee.Username)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if employee.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *employee.Password.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if employee.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (e CSEmployeeModel) Insert(csEmployee *CSEmployee) (string, error) {
	query := `
INSERT INTO cs_employee (first_name, last_name, username, password_hash, role, created_at) 
VALUES ($1, $2, $3, $4, $5, NOW()) 
RETURNING id`
	args := []interface{}{csEmployee.FirstName, csEmployee.LastName, csEmployee.Username, csEmployee.Password.hash, csEmployee.Role}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&csEmployee.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` ||
			err.Error() == `pq: повторяющееся значение ключа нарушает ограничение уникальности "users_email_key"`:
			return "", ErrDuplicateUsername
		default:
			return "", err
		}
	}
	return csEmployee.Username, nil
}

func (e CSEmployeeModel) GetByUsername(username string) (*CSEmployee, error) {
	query := `
SELECT id, first_name, last_name, username, password_hash, role, created_at, deleted_at, updated_at
FROM cs_employee
WHERE username = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var csEmployee CSEmployee

	err := e.DB.QueryRowContext(ctx, query, username).Scan(
		&csEmployee.ID,
		&csEmployee.FirstName,
		&csEmployee.LastName,
		&csEmployee.Username,
		&csEmployee.Password,
		&csEmployee.Role,
		&csEmployee.CreatedAt,
		&csEmployee.DeletedAt,
		&csEmployee.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &csEmployee, nil
}

func (e *CSEmployeeModel) GetAll() ([]*CSEmployee, error) {
	// Declare the SQL statement
	query := `
SELECT id, first_name, last_name, username, password_hash, role, created_at, deleted_at, updated_at 
FROM cs_employee`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []*CSEmployee{}

	// Iterate over the rows, scanning the result into a new employee struct and
	// appending it to the employees slice.
	for rows.Next() {
		var employee CSEmployee
		err := rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.Username,
			&employee.Role,
			&employee.CreatedAt,
			&employee.UpdatedAt,
			&employee.DeletedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, &employee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (e *CSEmployeeModel) Update(employee *CSEmployee) error {
	query := `
UPDATE cs_employee 
SET first_name = $1, last_name = $2, username = $3, 
password_hash = $4, role = $5, updated_at = NOW() 
WHERE id = $6`
	args := []interface{}{employee.FirstName, employee.LastName, employee.Username, employee.Password.hash, employee.Role, employee.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := e.DB.ExecContext(ctx, query, args...)
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

func (e *CSEmployeeModel) Delete(id int64) error {
	query := `
DELETE FROM cs_employee 
WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := e.DB.ExecContext(ctx, query, id)
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

func (e *CSEmployeeModel) SoftDelete(id int64) error {
	query := `
UPDATE cs_employee 
SET deleted_at = NOW() 
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := e.DB.ExecContext(ctx, query, id)
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