package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/saifoelloh/ranger/internal/constant"
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

func (r *UserRepository) FindBySSOID(ssoID string, ssoPlatform constant.SSOPlatform) (*model.User, error) {
	var user model.User
	var query string

	if ssoPlatform == constant.SSOPlatformApple {
		query = `
			SELECT id, first_name, last_name, email, phone_number, investor_type, sso_sign_option, apple_sso_id
			FROM "Users"
			WHERE apple_sso_id = $1 and sso_sign_option = $2`
	} else if ssoPlatform == constant.SSOPlatformGoogle {
		query = `
			SELECT id, first_name, last_name, email, phone_number, investor_type, sso_sign_option, google_sso_id
			FROM "Users"
			WHERE google_sso_id = $1 and sso_sign_option = $2`
	}
	err := r.db.Get(&user, query, ssoID, ssoPlatform)
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
