package repo

import (
	"context"
	"github.com/itzngga/Lara/entity"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"time"
)

var PremiumRepository *premiumRepository

type premiumRepository struct {
	DB *bun.DB
}

func NewPremiumRepository(DB *bun.DB) *premiumRepository {
	return &premiumRepository{DB: DB}
}

func (repository premiumRepository) GetAllPremiums(q string, page, limit int) ([]entity.PremiumEntity, int, error) {
	var list = make([]entity.PremiumEntity, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := repository.DB.NewSelect().Model(&list)
	if q != "" {
		parsedQ := "%" + q + "%"
		if repository.DB.Dialect().Name() == dialect.PG {
			query = query.Where("jid ILIKE ?", parsedQ)
		} else {
			query = query.Where("jid LIKE ?", parsedQ)
		}
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Scan(ctx); err != nil {
		return nil, 0, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	totalData, err := query.Count(ctx)
	if err != nil {
		return list, 0, err
	}

	return list, totalData, nil
}

func (repository premiumRepository) GetPremiumById(jid string) (entity.PremiumEntity, error) {
	var result entity.PremiumEntity

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := repository.DB.NewSelect().Model(&result).Where("jid = ?", jid).Scan(ctx); err != nil {
		return result, err
	}

	return result, nil
}

func (repository premiumRepository) InsertNewPremium(premiumEntity entity.PremiumEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewInsert().Model(&premiumEntity).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository premiumRepository) UpdatePremium(premiumEntity entity.PremiumEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewUpdate().Model(&premiumEntity).OmitZero().WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository premiumRepository) DeletePremium(jid []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewDelete().Model(&entity.PremiumEntity{}).Where("jid IN(?)", bun.In(jid)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository premiumRepository) IsValidPremium(jid string) (bool, error) {
	var ok bool
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := repository.DB.NewRaw("SELECT EXISTS(SELECT 1 FROM premium WHERE jid = ? LIMIT 1)", jid).Scan(ctx, &ok); err != nil {
		return false, err
	}

	return ok, nil
}
