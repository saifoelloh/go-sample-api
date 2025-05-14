package model

import (
	"database/sql"
	"time"
)

type Session struct {
	ID         string         `db:"id"`
	UserID     string         `db:"user_id"`
	Device     string         `db:"device"`
	MacAddress string         `db:"mac_address"`
	PublicKey  string         `db:"public_key"`
	Active     bool           `db:"active"`
	IP         sql.NullString `db:"ip"`
	UserAgent  sql.NullString `db:"user_agent"`
	Location   sql.NullString `db:"location"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}
