package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type AdminEntity struct {
	bun.BaseModel `bun:"table:admin"`

	Username  string    `bun:"username,pk"`
	Password  string    `bun:"password"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}
