package entity

import (
	"github.com/uptrace/bun"
)

type ReminderEntity struct {
	bun.BaseModel `bun:"table:reminder"`

	JID      string `bun:"jid,pk"`
	Duration int64  `bun:"duration"`
	Message  string `bun:"message"`
}
