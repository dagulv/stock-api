package domain

import (
	"time"

	"github.com/rs/xid"
)

type Credentials struct {
	UserId       xid.ID `json:"userId"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	OtpSecret    string `json:"otpSecret"`
	CredentialId string `json:"credentialId"`
	PublicKey    string `json:"publicKey"`
}

type Session struct {
	Id          xid.ID    `json:"id"`
	UserId      xid.ID    `json:"userId"`
	TimeExpired time.Time `json:"timeExpired"`
}

// TODO make tenantid a slice
type SessionUser struct {
	Id        xid.ID `json:"id"`
	TenantId  xid.ID `json:"tenantId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}