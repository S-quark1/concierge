package data

import (
	"context"
	"database/sql"
	"time"
)

type RegForm struct {
	ID          int64     `json:"id"`
	CompanyName string    `json:"company_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type RegFormModel struct {
	DB *sql.DB
}

func (r *RegFormModel) Insert(regForm *RegForm) error {
	query := `
INSERT INTO RegForm (company_name, email, phone_number)
VALUES ($1, $2, $3)
RETURNING id, created_at`
	args := []interface{}{regForm.CompanyName, regForm.Email, regForm.PhoneNumber}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&regForm.ID, &regForm.CreatedAt)
}

func (r *RegFormModel) GetByDateAsc() ([]*RegForm, error) {
	query := `
    SELECT id, company_name, email, phone_number, created_at
    FROM RegForm
    ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	forms := []*RegForm{}
	for rows.Next() {
		form := &RegForm{}
		err := rows.Scan(&form.ID, &form.CompanyName, &form.Email, &form.PhoneNumber, &form.CreatedAt)
		if err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	return forms, nil
}
