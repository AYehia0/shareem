// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Share struct {
	ID        uuid.UUID      `json:"id"`
	Url       string         `json:"url"`
	Title     string         `json:"title"`
	Note      sql.NullString `json:"note"`
	Ip        pgtype.Inet    `json:"ip"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
