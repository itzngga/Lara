package repo

import (
	"context"
	"github.com/itzngga/Lara/entity"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"time"
)

var AdminRepository *adminRepository

type adminRepository struct {
	DB *bun.DB
}

func NewAdminRepository(DB *bun.DB) *adminRepository {
	return &adminRepository{DB: DB}
}

func (repository adminRepository) GetAllAdmins(q string, page, limit int) ([]entity.AdminEntity, int, error) {
	var list = make([]entity.AdminEntity, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := repository.DB.NewSelect().Model(&list)
	if q != "" {
		parsedQ := "%" + q + "%"
		if repository.DB.Dialect().Name() == dialect.PG {
			query = query.Where("username ILIKE ?", parsedQ)
		} else {
			query = query.Where("username LIKE ?", parsedQ)
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

func (repository adminRepository) GetAdminById(id string) (entity.AdminEntity, error) {
	var result entity.AdminEntity

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := repository.DB.NewSelect().Model(&result).Where("username = ?", id).Scan(ctx); err != nil {
		return result, err
	}

	return result, nil
}

func (repository adminRepository) InsertNewAdmin(adminEntity entity.AdminEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewInsert().Model(&adminEntity).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository adminRepository) UpdateAdmin(adminEntity entity.AdminEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewUpdate().Model(&adminEntity).OmitZero().WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository adminRepository) DeleteAdmin(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := repository.DB.NewDelete().Model(&entity.AdminEntity{}).Where("username = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repository adminRepository) IsValidAdmin(id string) (bool, error) {
	var ok bool
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := repository.DB.NewRaw("SELECT EXISTS(SELECT 1 FROM admin WHERE username = ? LIMIT 1)", id).Scan(ctx, &ok); err != nil {
		return false, err
	}

	return ok, nil
}
