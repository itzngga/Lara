package conf

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/itzngga/Lara/entity"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"log"
	"os"
	"time"
)

var sqliteDB *bun.DB

func NewSqliteDB() *bun.DB {
	if sqliteDB != nil {
		return sqliteDB
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?_foreign_keys=on", os.Getenv("SQLITE_FILE")))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sqliteDB = bun.NewDB(sqldb, sqlitedialect.New())
	sqliteDB.SetMaxOpenConns(1)

	MigrateTables(&entity.ReminderEntity{}, &entity.WMEntity{})

	return sqliteDB
}

func MigrateTables(values ...interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	for _, value := range values {
		if _, err := sqliteDB.NewCreateTable().Model(value).IfNotExists().Exec(ctx); err != nil {
			log.Fatal(err)
		}
	}

	return
}
