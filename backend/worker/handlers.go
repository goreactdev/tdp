package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)


func (app *application) AddCollection(ctx context.Context, t *asynq.Task) error {

	var collection database.SBTCollection


	if err := json.Unmarshal(t.Payload(), &collection); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	insertedCollection, err := app.getCollection(collection)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get_collections: %v", err))
		return err
	}

	// for loop to next item index from 0 to next item index
	for i := int64(0); i < insertedCollection.NextItemIndex; i++ {
			payloadData := struct {
				CollectionAddress string `json:"collection_address"`
				ItemIndex 	   int64  `json:"item_index"`
			}{
				CollectionAddress: insertedCollection.FriendlyAddress,
				ItemIndex: i,
			}

			payload, err := json.Marshal(payloadData)

			if err != nil {
				app.logger.Error(err, nil)
				return err
			}

			runGetNfts := asynq.NewTask(database.TYPE_MIGRATE_NFT, payload)


			info, err := app.asynqClient.Enqueue(runGetNfts, asynq.MaxRetry(5),  asynq.ProcessIn(5*time.Second), asynq.Retention(10 * time.Minute), asynq.Queue(database.PRIORITY_URGENT))
			if err != nil {
				app.logger.Error(fmt.Errorf("error: %s", err), nil)
				return err
			}

			app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))
	}

	app.logger.Info(fmt.Sprintf("collection %s migrated", insertedCollection.FriendlyAddress))
	
	return nil
}


func (app *application) MigrateCollection(ctx context.Context, t *asynq.Task) error {

	var collectionAddr string


	if err := json.Unmarshal(t.Payload(), &collectionAddr); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// get collection from db

	collectionClient := nft.NewCollectionClient(app.tonLiteClient, address.MustParseAddr(collectionAddr))

	collectionData, err := collectionClient.GetCollectionData(ctx)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get_collection_data: %v", err))
		return err
	}

	// update collection in db
	err = app.sqlModels.Nfts.UpdateNextItemIndex(collectionAddr, collectionData.NextItemIndex.Int64())
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running update_collection: %v", err))
		return err
	}

	// for loop to next item index from 0 to next item index
	for i := int64(0); i < collectionData.NextItemIndex.Int64(); i++ {
			payloadData := struct {
				CollectionAddress string `json:"collection_address"`
				ItemIndex 	   int64  `json:"item_index"`
			}{
				CollectionAddress: collectionAddr,
				ItemIndex: i,
			}

			payload, err := json.Marshal(payloadData)

			if err != nil {
				app.logger.Error(err, nil)
				return err
			}

			runGetNfts := asynq.NewTask(database.TYPE_MIGRATE_NFT, payload)

			info, err := app.asynqClient.Enqueue(runGetNfts, asynq.MaxRetry(5),  asynq.ProcessIn(5*time.Second), asynq.Retention(10 * time.Minute), asynq.Queue(database.PRIORITY_URGENT))
			if err != nil {
				app.logger.Error(fmt.Errorf("error: %s", err), nil)
				return err
			}

			app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))
	}

	app.logger.Info(fmt.Sprintf("collection %s migrated", collectionAddr))
	
	return nil
}




func (app *application) MigrateNFT(ctx context.Context, t *asynq.Task) error {

	var payloadData struct {
		CollectionAddress string `json:"collection_address"`
		ItemIndex	 int64    `json:"item_index"`
	}

	if err := json.Unmarshal(t.Payload(), &payloadData); err != nil {
		app.logger.Warning(fmt.Sprintf("error unmarshalling payload: %v", err))
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	nft, err := app.getNFTByCollection(payloadData.CollectionAddress, payloadData.ItemIndex)
	if err != nil {
		return app.SkipError(err, t)
	}

	if nft == nil {
		return nil
	}

	payload := struct {
		UserAddr string `json:"user_address"`
		NftAddr string `json:"nft_address"`
	}{
		UserAddr: nft.FriendlyOwnerAddress,
		NftAddr: nft.FriendlyAddress,
	}

	outcomePayload, err := json.Marshal(payload)

	if err != nil {
		app.logger.Error(err, nil)
		return err
	}


	runGetRewardToAccount := asynq.NewTask(database.TYPE_ADD_REWARD_TO_ACCOUNT, outcomePayload)

	info, err := app.asynqClient.Enqueue(runGetRewardToAccount, asynq.TaskID(fmt.Sprintf(payload.NftAddr, payload.UserAddr)), asynq.MaxRetry(10),  asynq.ProcessIn(20*time.Second), asynq.Retention(24 * time.Hour), asynq.Queue(database.PRIORITY_URGENT))
	if err != nil {
		app.logger.Error(fmt.Errorf("error: %s", err), nil)
		return err
	}

	app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))	

	return nil
}


func (app *application) AddTgMessage(ctx context.Context, t *asynq.Task) error {

	var payloadData struct {
		UserId int `json:"user_id"`
		MessageId int `json:"message_id"`
		ChatId int64 `json:"chat_id"`
	}

	if err := json.Unmarshal(t.Payload(), &payloadData); err != nil {
		app.logger.Warning(fmt.Sprintf("error unmarshalling payload: %v", err))
		return err
	}

	// Get user by telegram id	

	tx, err := app.sqlModels.Rewards.DB.BeginTx(ctx, nil)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error starting transaction: %v", err))
		return err
	}

	id, err := app.sqlModels.Rewards.InsertTelegramMessage(tx, payloadData.UserId, payloadData.MessageId, payloadData.ChatId)
	if err != nil {
		tx.Rollback()
		app.logger.Warning(fmt.Sprintf("error running add tg message handler: %v", err))
		return app.SkipError(err, t)
	}

	// // // update user rating
	// err = app.sqlModels.Rewards.UpdateRating(tx, payloadData.UserId)
	// if err != nil {
	// 	tx.Rollback()
	// 	app.logger.Warning(fmt.Sprintf("error running update user rating handler: %v", err))
	// 	return err
	// }

	user, err := app.sqlModels.Users.GetByTelegramUserId(payloadData.UserId)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get user handler: %v", err))
		return err
	}


	// get prototype based user rating
	prototypesNFT, err := app.sqlModels.Nfts.GetPrototypesByRating(user.ID)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get user rating handler: %v", err))
		return err
	}

	for _, prototype := range prototypesNFT {
		// check if user has this nft
		id, err := app.sqlModels.Rewards.InsertStoredReward(user.FriendlyAddress, app.config.App.AdminCollectionAddress, prototype.Base64)
		if err != nil {
			app.logger.Warning(fmt.Sprintf("error running insert stored reward handler: %v", err))
			return err
		}
	     app.logger.Info(fmt.Sprintf("added stored reward with id %d", id))		
	}

	if err = tx.Commit(); err != nil {
		app.logger.Warning(fmt.Sprintf("error committing transaction: %v", err))
		return err
	}
	

	app.logger.Info(fmt.Sprintf("added tg message with id %d", id))

	return nil
}

func (app *application) RewardForLinkedAccounts(ctx context.Context, t *asynq.Task) error {

	var userId int64

	if err := json.Unmarshal(t.Payload(), &userId); err != nil {
		app.logger.Warning(fmt.Sprintf("error unmarshalling payload: %v", err))
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	nftMetadata, err := app.sqlModels.Nfts.GetNFTMetadataByID(app.config.App.AuthMetadataID)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get nft metadata by id handler: %v", err))
		return err
	}

	// get user by id
	user, err := app.sqlModels.Users.GetById(userId)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get user by id handler: %v", err))
		return err
	}

	// insert stored_reward
	id, err := app.sqlModels.Rewards.InsertStoredReward(user.FriendlyAddress, app.config.App.AdminCollectionAddress, nftMetadata.Base64)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running insert stored reward handler: %v", err))
		return err
	}

	// count stored rewards
	count, err := app.sqlModels.Rewards.CountStoredRewards()
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running count stored rewards handler: %v", err))
		return err
	}


	lastTime, err := app.sqlModels.Rewards.GetLastTimeStoredReward()
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get last time stored reward handler: %v", err))
		return err
	}

	// if count is equal to 20 or more or last time is more than 1 hour then send reward
	if count >= 20 || time.Now().Unix() - lastTime >= 3600 {

		runMintStoredRewards := asynq.NewTask(database.TYPE_MINT_STORED_REWARDS, nil)

		info, err := app.asynqClient.Enqueue(runMintStoredRewards, asynq.TaskID("MINT_STORED_REWARDS"), asynq.MaxRetry(10),  asynq.ProcessIn(5*time.Second), asynq.Retention(30 * time.Second), asynq.Queue(database.PRIORITY_URGENT))
		if err != nil {
			app.logger.Error(fmt.Errorf("error: %s", err), nil)
			return err
		}

		app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

		return nil
	}

	app.logger.Info(fmt.Sprintf("added stored reward with id %d", id))

	return nil
}

func (app *application) MintStoredRewards(ctx context.Context, t *asynq.Task) error {

	// // get all stored rewards
	rewards, err := app.sqlModels.Rewards.GetAllStoredRewards()
	if err != nil {
		app.logger.Error(err, nil)
		return err
	}

	// create map rewards by collection address
	rewardsByCollection := make(map[string][]*database.StoredReward)

	// must be like this: collectionAddress => rewards 
	// but in one reward is one collection address and it's can be duplicated
	for _, reward := range rewards {
		rewardsByCollection[reward.CollectionAddress] = append(rewardsByCollection[reward.CollectionAddress], reward)
	}

	app.logger.Info(fmt.Sprintf("rewards by collection: %v", rewardsByCollection))
	
	for collectionAddress, rewards := range rewardsByCollection {
		// get collection data
		collectionAddr := address.MustParseAddr(collectionAddress)
		collectionClient := nft.NewCollectionClient(app.tonLiteClient, collectionAddr)

		collectionData, err := collectionClient.GetCollectionData(context.Background())
		if err != nil {
			app.logger.Error(err, nil)
			return err
		}


		dict := cell.NewDict(64)

		var storedNftAddrAndUserId []struct {
			NftAddr string 
			UserAddr string
		}

		for i, reward := range rewards {
			offchainCon := &nft.ContentOffchain{
				URI: reward.Base64Metadata + "/meta.json",
			}

			con := cell.BeginCell().MustStoreStringSnake(offchainCon.URI).EndCell()
			// get collection data
		
			dict.Set(cell.BeginCell().MustStoreUInt(collectionData.NextItemIndex.Uint64()+uint64(i), 64).EndCell(), cell.BeginCell().
				MustStoreCoins(tlb.MustFromTON("0.04").NanoTON().Uint64()).
				MustStoreRef(
					cell.BeginCell().
						MustStoreAddr(address.MustParseAddr(reward.UserAddress)). // owner
						MustStoreRef(con).
						MustStoreAddr(address.MustParseAddr(reward.CollectionAddress)). // editor address: admin wallet
						EndCell()).
				EndCell())
			
			nftAddr, err := collectionClient.GetNFTAddressByIndex(ctx, big.NewInt(int64(collectionData.NextItemIndex.Uint64()+uint64(i))))
			if err != nil {
				app.logger.Error(err, nil)
				return err
			}

			storedNftAddrAndUserId = append(storedNftAddrAndUserId, struct {
				NftAddr string
				UserAddr string
			}{
				NftAddr: nftAddr.String(),
				UserAddr: reward.UserAddress,
			})

		}
		dataCell := cell.BeginCell().
			MustStoreUInt(2, 32).             // op code for mint batch
			MustStoreUInt(rand.Uint64(), 64). // query id
			MustStoreRef(dict.MustToCell()).
			EndCell()
			
		w := app.getWallet()

		mint := wallet.SimpleMessage(collectionAddr, tlb.MustFromTON(fmt.Sprint(0.06 * float64(len(rewards)))), dataCell)

		err = w.Send(ctx, mint, true)

		if err != nil {
			app.logger.Error(err, nil)
			return err
		}

		for _, object := range storedNftAddrAndUserId {
			payloadData := struct {
				UserAddr string `json:"user_address"`
				NftAddr string `json:"nft_address"`
			}{
				UserAddr: object.UserAddr,
				NftAddr: object.NftAddr,
			}

			payload, err := json.Marshal(payloadData)

			if err != nil {
				app.logger.Error(err, nil)
				return err
			}

			runGetRewardToAccount := asynq.NewTask(database.TYPE_ADD_REWARD_TO_ACCOUNT, payload)

			info, err := app.asynqClient.Enqueue(runGetRewardToAccount, asynq.TaskID(fmt.Sprintf(object.NftAddr, object.UserAddr)), asynq.MaxRetry(10),  asynq.ProcessIn(20*time.Second), asynq.Retention(24 * time.Hour), asynq.Queue(database.PRIORITY_URGENT))
			if err != nil {
				app.logger.Error(fmt.Errorf("error: %s", err), nil)
				return err
			}

			app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))	
		}
		

		app.logger.Info(fmt.Sprintf("minted %d nfts for collection %s", len(rewards), collectionAddress), nil)		

	}

	// mark reward as processed
	err = app.sqlModels.Rewards.MarkStoredRewardsAsProcessed()
	if err != nil {
		app.logger.Error(err, nil)
		return err
	}
	

	return nil
}
	
	


// 
func (app *application) SetReward(ctx context.Context, t *asynq.Task) error {

	var payloadData struct {
		UserAddr string `json:"user_address"`
		NftAddress string `json:"nft_address"`
	}

	if err := json.Unmarshal(t.Payload(), &payloadData); err != nil {
		app.logger.Warning(fmt.Sprintf("error unmarshalling payload: %v", err))
		return err
	}

	tx := app.sqlModels.Rewards.DB.MustBeginTx(ctx, nil)

	// check if nft already exists
	nft, err := app.sqlModels.Nfts.GetTokenByAddress(payloadData.NftAddress)
	if err != nil {
		tx.Rollback()
		app.logger.Error(err, nil)
		return err
	}

	// if nft does not exist - add it
	if nft == nil {
		// get nft metadata
		nft, err = app.insertNFTbyAddr(tx, payloadData.NftAddress)
		if err != nil {
			tx.Rollback()
			app.logger.Error(err, nil)
			return err
		}
	}	

	// get user by address
	user, err := app.sqlModels.Users.GetByFriendlyAddress(payloadData.UserAddr)
	if err != nil {
		tx.Rollback()
		app.logger.Error(err, nil)
		return err
	}

	// insert reward
	id, err := app.sqlModels.Rewards.Insert(tx, user.ID, nft.ID)
	if err != nil {
		tx.Rollback()
		app.logger.Error(err, nil)
		return err
	}

	app.logger.Info(fmt.Sprintf("added reward with id %d", id))

	// count awards that has user


	// update user rating
	err = app.sqlModels.Rewards.UpdateRatingByReward(tx, user.ID, nft.CreatedAt, nft.Weight)
	if err != nil {
		tx.Rollback()
		app.logger.Warning(fmt.Sprintf("error running update user rating handler: %v", err))
		return err
	}

	if err = tx.Commit(); err != nil {
		app.logger.Warning(fmt.Sprintf("error committing transaction: %v", err))
		return err
	}

	app.logger.Info(fmt.Sprintf("added nft with id %d", nft.ID))
		
	return nil

}