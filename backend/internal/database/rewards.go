package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type RewardModel struct {
	DB *sqlx.DB
}

type Reward struct {
	ID int64 `db:"id" json:"id"`
	UserID int64 `db:"user_id" json:"user_id"`
	SbtTokenID int64 `db:"sbt_token_id" json:"sbt_token_id"`
	CreatedAt int64 `db:"created_at" json:"created_at"`
	UpdatedAt int64 `db:"updated_at" json:"updated_at"`
	Weight int64 `json:"weight"`
	Version int64 `db:"version" json:"version"`
}

type TelegramMessage struct {
	ID int64 `db:"id" json:"id"`
	UserID int `db:"user_id" json:"user_id"`
	MessageID int `db:"message_id" json:"message_id"`
	ChatID int64 `db:"chat_id" json:"chat_id"`
	CreatedAt int64 `db:"created_at" json:"created_at"`
	UpdatedAt int64 `db:"updated_at" json:"updated_at"`
	Version int64 `db:"version" json:"version"`
}


type Merch struct {
	ID int64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Amount int64 `db:"amount" json:"amount"`
	UserID int64 `db:"user_id" json:"user_id"`
	Store string `db:"store" json:"store"`
	CreatedAt int64 `db:"created_at" json:"created_at"`
	UpdatedAt int64 `db:"updated_at" json:"updated_at"`	
}

// insert merch and substract rating from user take the rating from the user
func (m *RewardModel) InsertMerch(merch *Merch) error {

	sql := `INSERT INTO merch (name, amount, user_id, store, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.Exec(sql, merch.Name, merch.Amount, merch.UserID, merch.Store, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

// get merch by id
func (m *RewardModel) GetMerchByID(id int64) (*Merch, error) {
	var merch Merch

	err := m.DB.Get(&merch, "SELECT * FROM merch WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &merch, nil
}

// get all merches
func (m *RewardModel) GetAllMerch() ([]*Merch, error) {
	var merchs []*Merch

	err := m.DB.Select(&merchs, "SELECT * FROM merch")
	if err != nil {
		return nil, err
	}

	return merchs, nil
}



// get telegram messages by user id
func (m *RewardModel) CountTelegramMessagesByUserID(userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM tg_messages
	LEFT JOIN linked_accounts ON tg_messages.user_id = linked_accounts.telegram_user_id
	WHERE linked_accounts.user_id = $1
	`

	var count int64

	err := m.DB.Get(&count, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}



type StoredReward struct {
	ID int64 `db:"id" json:"id"`
	UserAddress string `db:"user_address" json:"user_address"`
	CollectionAddress string `db:"collection_address" json:"collection_address"`
	Base64Metadata string `db:"base64_metadata" json:"base64_metadata"`
	Processed bool `db:"processed" json:"processed"`
	ApprovedByUser bool `db:"approved_by_user" json:"approved_by_user"`	
	CreatedAt int64 `db:"created_at" json:"created_at"`
	UpdatedAt int64 `db:"updated_at" json:"updated_at"`
}


// insert reward

func (m *RewardModel) InsertStoredReward(userAddress string, collectionAddress string, base64Metadata string) (int64, error) {
	sql := `INSERT INTO stored_rewards (user_address, collection_address, base64_metadata, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id `

	var id int64

	err := m.DB.QueryRow(sql, userAddress, collectionAddress, base64Metadata, time.Now().Unix(), time.Now().Unix()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}




// count stored rewards

func (m *RewardModel) CountStoredRewards() (int64, error) {
	query := `SELECT COUNT(*) FROM stored_rewards WHERE processed = false AND approved_by_user = true`

	var count int64

	err := m.DB.QueryRow(query).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// get all achievements that are not processed and approved by user by user id
func (m *RewardModel) GetStoredRewardsByUserID(userAddr string, pagination *Pagination) ([]*StoredReward, error) {
	sql := `SELECT * FROM stored_rewards WHERE user_address = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var storedRewards []*StoredReward

	err := m.DB.Select(&storedRewards, sql, userAddr, pagination.End - pagination.Start, pagination.Start)
	if err != nil {
		return nil, err
	}

	return storedRewards, nil
}

// count stored achievements by user id

func (m *RewardModel) CountStoredRewardsByUserID(userAddr string) (int64, error) {
	sql := `SELECT COUNT(*) FROM stored_rewards WHERE user_address = $1`

	var count int64

	err := m.DB.QueryRow(sql, userAddr).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// get stored reward by id

func (m *RewardModel) GetStoredRewardByID(id int64) (*StoredReward, error) {
	sql := `SELECT * FROM stored_rewards WHERE id = $1`

	var storedReward StoredReward

	err := m.DB.Get(&storedReward, sql, id)
	if err != nil {
		return nil, err
	}

	return &storedReward, nil
}

// update approved by user

func (m *RewardModel) UpdateStoredRewardApprovedByUser(id int64, approved bool) error {
	sql := `UPDATE stored_rewards SET approved_by_user = $1 WHERE id = $2`

	_, err := m.DB.Exec(sql, approved, id)
	if err != nil {
		return err
	}

	return nil
}


// get last time of stored reward

func (m *RewardModel) GetLastTimeStoredReward() (int64, error) {
	query := `SELECT created_at FROM stored_rewards  
	WHERE processed = false AND approved_by_user = true
	ORDER BY created_at DESC LIMIT 1`

	var lastTime int64

	err := m.DB.QueryRow(query).Scan(&lastTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return lastTime, nil
}

// get all stored rewards

func (m *RewardModel) GetAllStoredRewards() ([]*StoredReward, error) {
	sql := `SELECT * FROM stored_rewards WHERE processed = false AND approved_by_user = true`

	var storedRewards []*StoredReward

	err := m.DB.Select(&storedRewards, sql)
	if err != nil {
		return nil, err
	}

	return storedRewards, nil
}

// mark stored rewards as processed

func (m *RewardModel) MarkStoredRewardsAsProcessed() error {
	sql := `UPDATE stored_rewards SET processed = true WHERE processed = false AND approved_by_user = true`

	_, err := m.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}



func (m *RewardModel) Insert(tx *sqlx.Tx, userId, sbtTokenId int64) (int64, error) {
	sql := `INSERT INTO rewards (user_id, sbt_token_id, created_at, updated_at, version) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int64

	err := tx.QueryRow(sql, userId, sbtTokenId, time.Now().Unix(), time.Now().Unix(), 1).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// get all rewards
// func (m *ActivitiesModel) GetAll(start, end int, filter string) ([]*Activity, error) {
// 	query := fmt.Sprintf(`SELECT * FROM activities %s ORDER BY id LIMIT $1 OFFSET $2`, filter)

// 	activities := []*Activity{}

// 	err := m.DB.Select(&activities, query, end, start)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return activities, nil

// }

// get all rewards
func (m *RewardModel) GetAll(start, end int, filter string) ([]*Reward, error) {
	query := fmt.Sprintf(`SELECT r.*, t.weight FROM rewards r
	LEFT JOIN sbt_tokens t ON t.id = r.sbt_token_id	
	%s ORDER BY id DESC LIMIT $1 OFFSET $2`, filter)

	rewards := []*Reward{}

	err := m.DB.Select(&rewards, query, end, start)
	if err != nil {
		return nil, err
	}

	return rewards, nil

}

// get by id
func (m *RewardModel) GetById(id int64) (*Reward, error) {
	query := `SELECT * FROM rewards WHERE id = $1`

	reward := &Reward{}

	err := m.DB.Get(reward, query, id)
	if err != nil {
		return nil, err
	}

	return reward, nil
}



func (m *RewardModel) InsertTelegramMessage(tx *sql.Tx, userId, messageId int, chatId int64) (int64, error) {
	sql := `INSERT INTO tg_messages (user_id, message_id, chat_id, created_at, updated_at, version) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int64

	err := tx.QueryRow(sql, userId, messageId, chatId, time.Now().Unix(), time.Now().Unix(), 1).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// update user rating by fetching count and +1
func (m *RewardModel) UpdateRating(tx *sql.Tx, telegramUserId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE users SET rating = rating + 1
		WHERE id = (
			SELECT user_id FROM linked_accounts
			WHERE linked_accounts.telegram_user_id = $1
		)
		`

	_, err := tx.ExecContext(ctx, query, telegramUserId)

	if err != nil {
		return err
	}

	return nil
}

func (m *RewardModel) UpdateRatingByReward(tx *sqlx.Tx, userId, lastAwardsAt, weight int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE users
		SET rating = rating + $1, awards_count = awards_count + 1, last_award_at = $2
		WHERE id = $3;
			`

	_, err := tx.ExecContext(ctx, query, weight, lastAwardsAt, userId)

	if err != nil {
		return err
	}
	

	return nil
}


// count rewards by user id
func (m *RewardModel) CountByUserId(userId int64) (int64, error) {
	query := `SELECT COUNT(*) FROM rewards WHERE user_id = $1`

	var count int64

	err := m.DB.Get(&count, query, userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}




func (m *RewardModel) GetByUserID(userID int64) (*Reward, error) {
	sql := `SELECT * FROM rewards WHERE user_id = $1`

	var reward Reward

	err := m.DB.Get(&reward, sql, userID)
	if err != nil {
		return nil, err
	}

	return &reward, nil
}

func (m *RewardModel) GetByUserIDAndSbtTokenID(userID int64, sbtTokenID int64) (*Reward, error) {
	sql := `SELECT * FROM rewards WHERE user_id = $1 AND sbt_token_id = $2`

	var reward Reward

	err := m.DB.Get(&reward, sql, userID, sbtTokenID)
	if err != nil {
		return nil, err
	}

	return &reward, nil
}

func (m *RewardModel) GetBySbtTokenID(sbtTokenID int64) (*Reward, error) {
	sql := `SELECT * FROM rewards WHERE sbt_token_id = $1`

	var reward Reward

	err := m.DB.Get(&reward, sql, sbtTokenID)
	if err != nil {
		return nil, err
	}

	return &reward, nil
}

func (m *RewardModel) Update(reward Reward) error {
	sql := `UPDATE rewards SET user_id = $1, sbt_token_id = $2, created_at = $3, updated_at = $4, version = $5 WHERE id = $6`

	_, err := m.DB.Exec(sql, reward.UserID, reward.SbtTokenID, reward.CreatedAt, reward.UpdatedAt, reward.Version, reward.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *RewardModel) Delete(id int64) error {
	sql := `DELETE FROM rewards WHERE id = $1`

	_, err := m.DB.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}


