package repo

import (
	"context"
	"github.com/itzngga/Lara/conf"
	"github.com/itzngga/Lara/entity"
	"time"
)

type reminderRepository struct {
}

func (repository reminderRepository) InsertOrUpdateReminder(reminderEntity entity.ReminderEntity) error {
	var db = conf.NewSqliteDB("")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := db.NewRaw("INSERT INTO reminder VALUES(?, ?, ?) ON CONFLICT(jid) DO UPDATE SET duration = ?, message = ? WHERE reminder.jid = ?",
		reminderEntity.JID, reminderEntity.Duration, reminderEntity.Message, reminderEntity.Duration, reminderEntity.Message, reminderEntity.JID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository reminderRepository) GetAllReminders() ([]entity.ReminderEntity, error) {
	var list = make([]entity.ReminderEntity, 0)
	var db = conf.NewSqliteDB("")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err := db.NewSelect().Model(&list).Scan(ctx)
	if err != nil {
		return list, err
	}

	return list, nil
}
