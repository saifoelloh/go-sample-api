package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/saifoelloh/ranger/internal/model"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(session *model.Session) error {
	query := `
		INSERT INTO "Sessions"
		(id, user_id, client_version, device, mac_address, public_key, active, ip, user_agent)
		VALUES (:id, :user_id, :client_version, :device, :mac_address, :public_key, :active, :ip, :user_agent)
	`
	_, err := r.db.NamedExec(query, session)
	if err != nil {
		return errors.InternalServerError(
			errors.WithScope("SessionRepository"),
			errors.WithLocation("CreateSession"),
			errors.WithDetail(err.Error()),
			errors.WithErrorCode("session/create"),
		)
	}

	return nil
}

func (r *SessionRepository) DeactivateSessionsByUserID(userID string) error {
	query := `UPDATE "Sessions" SET active = false WHERE user_id = $1 AND active = true`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return errors.InternalServerError(
			errors.WithScope("SessionRepository"),
			errors.WithLocation("DeactivateSessionsByUserID"),
			errors.WithDetail(err.Error()),
			errors.WithErrorCode("session/deactivation-failed"),
		)
	}

	return nil
}
