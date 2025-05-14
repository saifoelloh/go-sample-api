package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                    string         `db:"id"`
	FirstName             string         `db:"first_name"`
	LastName              string         `db:"last_name"`
	Email                 sql.NullString `db:"email"`
	EmailHash             sql.NullString `db:"email_hash"`
	EmailVerificationCode string         `db:"email_verification_code"`
	EmailVerified         bool           `db:"email_verified"`
	Password              sql.NullString `db:"password"`
	PhoneNumber           sql.NullString `db:"phone_number"`
	PhoneNumberHash       sql.NullString `db:"phone_number_hash"`
	PhoneNumberVerified   bool           `db:"phone_number_verified"`
	LastChangePassword    time.Time      `db:"last_change_password"`
	InvestorType          sql.NullString `db:"investor_type"`
	IsDeleted             bool           `db:"is_deleted"`
	PhotoProfile          sql.NullString `db:"photo_profile"`
	Role                  string         `db:"role"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
}
