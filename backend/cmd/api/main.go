package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/ton-developer-program/internal/database"
	"github.com/ton-developer-program/internal/leveledlog"
	"github.com/ton-developer-program/internal/smtp"
	"github.com/ton-developer-program/internal/tonconnect"
	"github.com/ton-developer-program/internal/version"
	"github.com/ton-developer-program/util"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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
	config            util.Config
	sqlModels         database.Models
	logger            *leveledlog.Logger
	mailer            *smtp.Mailer
	wg                sync.WaitGroup
	asynqClient       *asynq.Client
	asynqScheduler    *asynq.Scheduler
	tonLiteClient     *ton.APIClient
	githubOauthConfig *oauth2.Config
	oauthStateString  string
}

func run(logger *leveledlog.Logger) error {

	var err error
	networks["-239"], err = liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatal(err)
	}
	networks["-3"], err = liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		log.Fatal(err)
	}

	// initialize config

	cfg, err := util.LoadConfig()
	if err != nil {
		return err
	}

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.Database.Dsn, cfg.Database.Automigrate)
	if err != nil {
		return err
	}
	defer db.Close()

	// initialize mailer

	mailer := smtp.NewMailer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From)

	// initialize task manager

	asynqClient := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			PoolSize: cfg.Redis.PoolSize,
		},
	)

	// initialize scheduler

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}

	asynqScheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			PoolSize: cfg.Redis.PoolSize,
		},
		&asynq.SchedulerOpts{
			Location: loc,
		},
	)

	githubOauthConfig := &oauth2.Config{
		RedirectURL:  cfg.Auth.GithubRedirectUrl,
		ClientID:     cfg.Auth.GithubClientId,
		ClientSecret: cfg.Auth.GithubClientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// instantiate application

	connectionPool := liteclient.NewConnectionPool()

	tonLiteClient, err := tonconnect.NewTonConnection(connectionPool, cfg)
	if err != nil {
		logger.Error(fmt.Errorf("error connecting to lite servers: %v", err), nil)
		return err
	}

	go func() {

		h := asynqmon.New(asynqmon.Options{
			RootPath:     "/monitoring", // RootPath specifies the root for asynqmon app
			RedisConnOpt: asynq.RedisClientOpt{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password},
		})
		

		http.Handle(h.RootPath()+"/", basicAuth(h, cfg.App.BasicUsername, cfg.App.BasicPassword, "Please enter your username and password"))

		fmt.Print("Starting asynqmon server.. \n")

		logger.Info(fmt.Sprintf("%s:8080%s", cfg.App.BaseUrl, h.RootPath()))

		// Go to http://localhost:8080/monitoring to see asynqmon homepage.
		http.ListenAndServe(":8080", nil)
	}()

	app := &application{
		config:            cfg,
		sqlModels:         database.NewModels(db.DB),
		logger:            logger,
		mailer:            mailer,
		asynqClient:       asynqClient,
		tonLiteClient:     tonLiteClient,
		asynqScheduler:    asynqScheduler,
		githubOauthConfig: githubOauthConfig,
		oauthStateString:  "erbEKBi3w4oirewbikjewrbuio2wkwsvjeierorbbre",
	}

	return app.serveHTTP()
}
