package entity

import "github.com/uptrace/bun"

type WMEntity struct {
	bun.BaseModel `bun:"table:sticker_wm"`

	JID              string `bun:"jid,pk"`
	StickerName      string `bun:"sticker_name"`
	StickerPublisher string `bun:"sticker_publisher"`
}
