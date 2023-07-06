package repo

import (
	"context"
	"github.com/itzngga/Lara/entity"
	"github.com/uptrace/bun"
	"time"
)

var ReminderRepository *reminderRepository

type reminderRepository struct {
	DB *bun.DB
}

func NewReminderRepository(DB *bun.DB) *reminderRepository {
	return &reminderRepository{DB: DB}
}

func (repository reminderRepository) PutReminder(reminderEntity entity.ReminderEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewRaw("INSERT INTO reminder VALUES(?, ?, ?) ON CONFLICT(jid) DO UPDATE SET duration = ?, message = ? WHERE reminder.jid = ?",
		reminderEntity.JID, reminderEntity.Duration, reminderEntity.Message, reminderEntity.Duration, reminderEntity.Message, reminderEntity.JID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository reminderRepository) GetAllReminders() ([]entity.ReminderEntity, error) {
	var list = make([]entity.ReminderEntity, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err := repository.DB.NewSelect().Model(&list).Scan(ctx)
	if err != nil {
		return list, err
	}

	return list, nil
}
