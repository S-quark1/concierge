package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Service     ServiceModel
	Price       PriceModel
	Company     CompanyModel
	User        UserModel
	RegForm     RegFormModel
	Token       TokenModel
	Permissions PermissionModel
	Request     RequestModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Service:     ServiceModel{DB: db},
		Price:       PriceModel{DB: db},
		Company:     CompanyModel{DB: db},
		User:        UserModel{DB: db},
		RegForm:     RegFormModel{DB: db},
		Token:       TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Request:     RequestModel{DB: db},
	}
}
