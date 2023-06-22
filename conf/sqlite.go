package conf

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/itzngga/Lara/entity"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"time"
)

var sqliteDB *bun.DB

func NewSqliteDB(dbName string) *bun.DB {
	if sqliteDB != nil {
		return sqliteDB
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?_foreign_keys=on", dbName))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sqliteDB = bun.NewDB(sqldb, sqlitedialect.New())
	sqliteDB.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if _, err := sqliteDB.NewCreateTable().Model(&entity.ReminderEntity{}).IfNotExists().Exec(ctx); err != nil {
		fmt.Println(err)
		return sqliteDB
	}

	return sqliteDB
}
