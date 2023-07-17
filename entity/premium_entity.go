package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type PremiumEntity struct {
	bun.BaseModel `bun:"table:premium"`

	JID       string    `bun:"jid,pk"`
	CreatedAt time.Time `bun:"created_at"`
	ExpiredAt time.Time `bun:"expired_at"`
}
