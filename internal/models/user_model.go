package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/xid"
)

type User struct {
	Id          xid.ID             `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Active      bool               `json:"active"`
	TimeCreated pgtype.Timestamptz `json:"timeCreated"`
	TimeUpdated pgtype.Timestamptz `json:"timeUpdated"`
}
