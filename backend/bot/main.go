package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/ton-developer-program/internal/leveledlog"
	"github.com/ton-developer-program/internal/version"
	"github.com/ton-developer-program/util"
)

func main() {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)

	err := run(logger)
	if err != nil {
		trace := debug.Stack()
		logger.Fatal(err, trace)
	}
}

type application struct {
	config util.Config
	sqlModels     database.Models
	logger *leveledlog.Logger
	wg     sync.WaitGroup
	asynqClient *asynq.Client
	bot *tgbotapi.BotAPI
}

func run(logger *leveledlog.Logger) error {

	config, err := util.LoadConfig()
	if err != nil {
		return err
	}
//dd
	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()
	
	asynqClient := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
		},
	)

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(config.Database.Dsn, false)
	if err != nil {
		return err
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(config.Auth.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	app := &application{
		config: config,
		sqlModels: database.NewModels(db.DB),    
		logger: logger,
		asynqClient: asynqClient,
		bot:    bot,
	}

	return app.startBot()

}
