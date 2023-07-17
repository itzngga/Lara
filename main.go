package main

import (
	"github.com/itzngga/Lara/conf"
	"github.com/itzngga/Lara/repo"
	_ "github.com/itzngga/Lara/src/cmd"
	"log"
	"time"

	"github.com/itzngga/Roxy/core"
	"github.com/itzngga/Roxy/options"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	"os"
	"os/signal"
	_ "time/tzdata"

	"syscall"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
}

func main() {
	bunDB := conf.NewSqliteDB()

	reminderRepository := repo.NewReminderRepository(bunDB)
	premiumRepository := repo.NewPremiumRepository(bunDB)
	adminRepository := repo.NewAdminRepository(bunDB)
	wmRepository := repo.NewWMRepository(bunDB)

	repo.ReminderRepository = reminderRepository
	repo.PremiumRepository = premiumRepository
	repo.AdminRepository = adminRepository
	repo.WMRepository = wmRepository

	roxyOptions := options.Options{
		StoreMode:                   "sqlite",
		LogLevel:                    "INFO",
		WithSqlDB:                   bunDB.DB,
		WithCommandLog:              true,
		AllowFromGroup:              true,
		AllowFromPrivate:            true,
		CommandSuggestion:           true,
		CommandResponseCacheTimeout: time.Minute * 15,
		SendMessageTimeout:          time.Second * 15,
	}

	app, err := core.NewGoRoxyBase(&roxyOptions)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
