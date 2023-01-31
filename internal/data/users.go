package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"github.com/concierge/service/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)
var AnonymousUser = &User{}

type User struct {
	ID          int64          `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Email       string         `json:"email"`
	Username    string         `json:"username"`
	Password    password       `json:"password_hash"`
	Activated   bool           `json:"activated"`
	UserType    string         `json:"user_type"`
	Preferences string         `json:"preferences"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   sql.NullString `json:"deleted_at"`
	UpdatedAt   sql.NullString `json:"updated_at"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserModel struct {
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

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "firstname", "must be provided")
	v.Check(len(user.FirstName) <= 500, "firstname", "must not be more than 500 bytes long")

	v.Check(user.LastName != "", "lastname", "must be provided")
	v.Check(len(user.LastName) <= 500, "lastname", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (u UserModel) Insert(user *User) error {
	query := `
INSERT INTO users (first_name, last_name, email, username, password_hash, activated, user_type, preferences, created_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()) 
RETURNING id`
	args := []interface{}{user.FirstName, user.LastName, user.Email, user.Username, user.Password.hash, user.Activated, user.UserType, user.Preferences}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
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

func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `
SELECT id, first_name, last_name, email, username, password_hash, activated, user_type, preferences, created_at, deleted_at, updated_at
FROM users
WHERE email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var csEmployee User

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&csEmployee.ID,
		&csEmployee.FirstName,
		&csEmployee.LastName,
		&csEmployee.Email,
		&csEmployee.Username,
		&csEmployee.Password.hash,
		&csEmployee.Activated,
		&csEmployee.UserType,
		&csEmployee.Preferences,
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

func (u UserModel) GetByUsername(username string) (*User, error) {
	query := `
SELECT id, first_name, last_name, email, username, password_hash, activated, user_type, preferences, created_at, deleted_at, updated_at
FROM users
WHERE username = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var csEmployee User

	err := u.DB.QueryRowContext(ctx, query, username).Scan(
		&csEmployee.ID,
		&csEmployee.FirstName,
		&csEmployee.LastName,
		&csEmployee.Email,
		&csEmployee.Username,
		&csEmployee.Password.hash,
		&csEmployee.Activated,
		&csEmployee.UserType,
		&csEmployee.Preferences,
		&csEmployee.CreatedAt,
		&csEmployee.DeletedAt,
		&csEmployee.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			print("bruh\n")
			return nil, err
		}
	}
	return &csEmployee, nil
}

func (u UserModel) GetAll() ([]*User, error) {
	// Declare the SQL statement
	query := `
SELECT id, first_name, last_name, email, username, password_hash, activated, user_type, preferences, created_at, deleted_at, updated_at
FROM users`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []*User{}

	// Iterate over the rows, scanning the result into a new employee struct and
	// appending it to the employees slice.
	for rows.Next() {
		var csEmployee User
		err := rows.Scan(
			&csEmployee.ID,
			&csEmployee.FirstName,
			&csEmployee.LastName,
			&csEmployee.Email,
			&csEmployee.Username,
			&csEmployee.Password.hash,
			&csEmployee.Activated,
			&csEmployee.UserType,
			&csEmployee.Preferences,
			&csEmployee.CreatedAt,
			&csEmployee.DeletedAt,
			&csEmployee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, &csEmployee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (u UserModel) Update(user *User) error {
	query := `
UPDATE users
SET first_name = $1, last_name = $2, email = $3, username = $4, password_hash = $5, activated = $6, user_type = $7, preferences = $8, updated_at = NOW()
WHERE id = $9 AND deleted_at IS NULL
RETURNING updated_at`

	args := []interface{}{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Username,
		user.Password.hash,
		user.Activated,
		user.UserType,
		user.Preferences,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
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

func (u UserModel) Delete(id int64) error {
	query := `
DELETE FROM users
WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, query, id)
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

func (u UserModel) SoftDelete(id int64) error {
	query := `
UPDATE users
SET deleted_at = NOW() 
WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, query, id)
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

func (u UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	query := `
SELECT users.id, first_name, last_name, email, username, password_hash, activated, user_type, preferences, created_at, deleted_at, updated_at
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`
	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password.hash,
		&user.Activated,
		&user.UserType,
		&user.Preferences,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}
