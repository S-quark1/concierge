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

type Employee struct {
	ID             int64     `json:"id"`
	FirstName      string    `json:"firstname"`
	LastName       string    `json:"lastname"`
	Username       string    `json:"username"`
	Password       password  `json:"-"`
	DateOfCreation time.Time `json:"-"`
	Role           int32     `json:"role"`
	Version        int32     `json:"version"`
}

type EmployeeModel struct {
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
func ValidateEmployee(v *validator.Validator, employee *Employee) {
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

func (e EmployeeModel) Insert(employee *Employee) error {
	query := `
INSERT INTO cs_employee (firstname, lastname, username, password_hash, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, dateOfCreation, version`
	args := []interface{}{employee.FirstName, employee.LastName, employee.Username, employee.Password.hash, employee.Role}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(&employee.ID, &employee.DateOfCreation, &employee.Version)
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

func (e EmployeeModel) GetByUsername(username string) (*Employee, error) {
	query := `
SELECT id, firstName, lastName, username, password_hash, role, dateOfCreation, version
FROM cs_employee
WHERE username = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var employee Employee

	err := e.DB.QueryRowContext(ctx, query, username).Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&employee.Username,
		&employee.Password.hash,
		&employee.Role,
		&employee.DateOfCreation,
		&employee.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &employee, nil
}
