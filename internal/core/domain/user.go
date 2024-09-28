package domain

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
)

// TODO make tenantid a slice
type User struct {
	Id          xid.ID    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email" validate:"required,email"`
	Active      bool      `json:"active"`
	TimeCreated time.Time `json:"timeCreated"`
	TimeUpdated time.Time `json:"timeUpdated"`
}

func (u *User) EncodeToStream(s *jsoniter.Stream) {
	s.WriteObjectField("id")
	s.WriteString(u.Id.String())

	s.WriteMore()
	s.WriteObjectField("firstName")
	s.WriteString(u.FirstName)

	s.WriteMore()
	s.WriteObjectField("lastName")
	s.WriteString(u.LastName)

	s.WriteMore()
	s.WriteObjectField("email")
	s.WriteString(u.Email)

	s.WriteMore()
	s.WriteObjectField("active")
	s.WriteBool(u.Active)

	s.WriteMore()
	s.WriteObjectField("timeCreated")
	s.WriteString(u.TimeCreated.Format(time.RFC3339))

	s.WriteMore()
	s.WriteObjectField("timeUpdated")
	s.WriteString(u.TimeUpdated.Format(time.RFC3339))
}
