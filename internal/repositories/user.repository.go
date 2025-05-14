package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/saifoelloh/ranger/internal/model"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, first_name, last_name, email, phone_number, investor_type, password
		FROM "Users"
		WHERE email_hash = $1`
	err := r.db.Get(&user, query, email)

	if err != nil {
		return nil, errors.NotFound(
			errors.WithScope("UserRepository"),
			errors.WithLocation("FindByEmail"),
			errors.WithMessage("user not found"),
			errors.WithErrorCode("user/not-found"),
		)
	}

	return &user, nil
}

func (r *UserRepository) FindBySSOID(ssoID string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, first_name, last_name, email, phone_number, investor_type, password
		FROM "Users"
		WHERE sso_id = $1`
	err := r.db.Get(&user, query, ssoID)
	if err != nil {
		return nil, errors.NotFound(
			errors.WithScope("UserRepository"),
			errors.WithLocation("FindBySSOID"),
			errors.WithMessage("user not found"),
			errors.WithErrorCode("user/not-found"),
		)
	}

	return &user, nil
}
