package database

import (
	"errors"
	"time"

	"github.com/ton-developer-program/assets"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

const defaultTimeout = 10 * time.Second

type DB struct {
	*sqlx.DB
}

const (
	TYPE_LISTENING_TRANSACTION = "master:listening_transaction"
	TYPE_MIGRATE_COLLECTION  = "master:migrate_collection"
	TYPE_ADD_COLLECTION  = "master:add_collection"
	
	TYPE_ADD_TG_MESSAGE  = "master:add_tg_message"
	TYPE_REWARD_FOR_LINKED_ACCOUNT  = "master:reward_for_linked_account"
	TYPE_ADD_REWARD_TO_ACCOUNT  = "master:add_reward_to_account"
	TYPE_MINT_STORED_REWARDS  = "master:mint_stored_rewards"

	TYPE_MIGRATE_NFT = "master:migrate_nft"
)

const (
	PRIORITY_CRITICAL = "critical"
	PRIORITY_URGENT = "urgent"
	PRIORITY_NORMAL = "normal"
	PRIORITY_LOW = "low"	
)

type Pagination struct {
	Start int
	End   int
}

func New(dsn string, automigrate bool) (*DB, error) {
	db, err := sqlx.Connect("postgres", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if automigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}

	return &DB{db}, nil
}
