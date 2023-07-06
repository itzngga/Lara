package repo

import (
	"context"
	"github.com/itzngga/Lara/entity"
	"github.com/uptrace/bun"
	"time"
)

var WMRepository *wmRepository

type wmRepository struct {
	DB *bun.DB
}

func NewWMRepository(DB *bun.DB) *wmRepository {
	return &wmRepository{DB: DB}
}

func (repository wmRepository) GetWMByJid(jid string) (entity.WMEntity, error) {
	var result entity.WMEntity

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := repository.DB.NewSelect().Model(&result).Where("jid = ?", jid).Scan(ctx); err != nil {
		return result, err
	}

	return result, nil
}

func (repository wmRepository) PutWM(wmEntity entity.WMEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewRaw("INSERT INTO sticker_wm VALUES(?, ?, ?) ON CONFLICT(jid) DO UPDATE SET sticker_name = ?, sticker_publisher = ? WHERE sticker_wm.jid = ?",
		wmEntity.JID, wmEntity.StickerName, wmEntity.StickerPublisher, wmEntity.StickerName, wmEntity.StickerPublisher, wmEntity.JID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
