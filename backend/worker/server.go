package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/ton-developer-program/internal/leveledlog"
	"github.com/ton-developer-program/internal/tonconnect"
	"github.com/ton-developer-program/internal/version"
	"github.com/ton-developer-program/util"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)



type application struct {
	config util.Config
	logger *leveledlog.Logger
	tonLiteClient *ton.APIClient
	asynqClient *asynq.Client
	sqlModels database.Models
}

func main() {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)

	err := run(logger)
	if err != nil {
		trace := debug.Stack()
		logger.Fatal(err, trace)
	}
}

func run(logger *leveledlog.Logger) error {

	cfg, err := util.LoadConfig()
	if err != nil {
		return err
	}

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	// Display version and exit.
	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

    // Set up connection using the tonLiteClient config


	connectionPool := liteclient.NewConnectionPool()
	
    tonLiteClient, err := tonconnect.NewTonConnection(connectionPool, cfg)
    if err != nil {
        logger.Error(fmt.Errorf("error connecting to lite servers: %v", err), nil)
        return err
    }

	asynqClient := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			PoolSize: cfg.Redis.PoolSize,
		},
	)


	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			PoolSize: cfg.Redis.PoolSize,
		},
		asynq.Config{Concurrency: cfg.Ton.MaxConcurrentTask, RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			return 5 * time.Second
		},
			Queues: map[string]int{
				database.PRIORITY_CRITICAL: 4,
				database.PRIORITY_URGENT:   3,
				database.PRIORITY_NORMAL:  2,
				database.PRIORITY_LOW:      1,
			},
		},
	)

	db, err := database.New(cfg.Database.Dsn, cfg.Database.Automigrate)
	if err != nil {
		return err
	}
	defer db.Close()


	app := &application{
		config: cfg,
		logger: logger,
		tonLiteClient: tonLiteClient,
		asynqClient: asynqClient,
		sqlModels: database.NewModels(db.DB),
	}

	stop := make(chan struct{})
    
    go func() {
        for {
            if !app.runGoroutine(stop) {
                // If the goroutine stops and the function returns false,
                // we wait for a bit before starting it again.
                time.Sleep(time.Second * 10)
            } else {
                // If the function returns true, we stop the loop.
                break
            }
        }
    }()
		

	if err := srv.Run(app.routes()); err != nil {
		log.Fatal(err)
	}

	return nil

}



func (app *application) routes() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(database.TYPE_ADD_COLLECTION, app.AddCollection)
	mux.HandleFunc(database.TYPE_MIGRATE_COLLECTION, app.MigrateCollection)

	mux.HandleFunc(database.TYPE_MIGRATE_NFT, app.MigrateNFT)
	mux.HandleFunc(database.TYPE_ADD_TG_MESSAGE, app.AddTgMessage)

	mux.HandleFunc(database.TYPE_REWARD_FOR_LINKED_ACCOUNT, app.RewardForLinkedAccounts)

	mux.HandleFunc(database.TYPE_MINT_STORED_REWARDS, app.MintStoredRewards)
	
	mux.HandleFunc(database.TYPE_ADD_REWARD_TO_ACCOUNT, app.SetReward)

	return mux
}
