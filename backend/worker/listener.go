package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

// func to get storage map key
func getShardID(shard *ton.BlockIDExt) string {
	return fmt.Sprintf("%d|%d", shard.Workchain, shard.Shard)
}

func getNotSeenShards(ctx context.Context, api *ton.APIClient, shard *ton.BlockIDExt, shardLastSeqno map[string]uint32) (ret []*ton.BlockIDExt, err error) {
	if no, ok := shardLastSeqno[getShardID(shard)]; ok && no == shard.SeqNo {
		return nil, nil
	}

	b, err := api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, fmt.Errorf("get block data: %w", err)
	}

	parents, err := b.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, fmt.Errorf("get parent blocks (%d:%x:%d): %w", shard.Workchain, uint64(shard.Shard), shard.Shard, err)
	}

	for _, parent := range parents {
		ext, err := getNotSeenShards(ctx, api, parent, shardLastSeqno)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ext...)
	}

	ret = append(ret, shard)
	return ret, nil
}

func (app *application) RunListeningTransactions(ctx context.Context) error {

	master, err := app.tonLiteClient.GetMasterchainInfo(ctx)
	if err != nil {
		app.logger.Error(errors.New("get masterchain info:"+err.Error()), nil)
		return err
	}

	// bound all requests to single lite server for consistency,
	// if it will go down, another lite server will be used
	context := app.tonLiteClient.Client().StickyContext(ctx)

	// storage for last seen shard seqno
	shardLastSeqno := map[string]uint32{}

	// getting information about other work-chains and shards of first master block
	// to init storage of last seen shard seq numbers
	firstShards, err := app.tonLiteClient.GetBlockShardsInfo(context, master)
	if err != nil {
		app.logger.Error(errors.New("get shards info:"+err.Error()), nil)
		return err
	}

	for _, shard := range firstShards {
		shardLastSeqno[getShardID(shard)] = shard.SeqNo
	}


	for {

		app.logger.Info("scanning new master block...")

		pagination := &database.Pagination{
			Start:  0,
			End: 100,
		}
				// get collections
		collections, err := app.sqlModels.Nfts.GetAllCollections(pagination)
		if err != nil {
			app.logger.Error(errors.New("get collections:"+err.Error()), nil)
			return err
		}	



		// getting information about other work-chains and shards of master block
		currentShards, err := app.tonLiteClient.GetBlockShardsInfo(context, master)
		if err != nil {
			app.logger.Error(errors.New("get shards info:"+err.Error()), nil)
			return err
		}

		// shards in master block may have holes, e.g. shard seqno 2756461, then 2756463, and no 2756462 in master chain
		// thus we need to scan a bit back in case of discovering a hole, till last seen, to fill the misses.
		var newShards []*ton.BlockIDExt
		for _, shard := range currentShards {
			notSeen, err := getNotSeenShards(context, app.tonLiteClient, shard, shardLastSeqno)
			if err != nil {
				app.logger.Error(errors.New("get not seen shards:"+err.Error()), nil)
				return err
			}
			shardLastSeqno[getShardID(shard)] = shard.SeqNo
			newShards = append(newShards, notSeen...)
		}

		var txList []*tlb.Transaction

		// for each shard block getting transactions
		for _, shard := range newShards {
			app.logger.Info("scanning shard block", "shard", shard.Shard, "seqno", shard.SeqNo)

			var fetchedIDs []ton.TransactionShortInfo
			var after *ton.TransactionID3
			var more = true

			// load all transactions in batches with 100 transactions in each while exists
			for more {
				fetchedIDs, more, err = app.tonLiteClient.WaitForBlock(master.SeqNo).GetBlockTransactionsV2(context, shard, 100, after)
				if err != nil {
					app.logger.Error(errors.New("get tx ids:"+err.Error()), nil)
					return err
				}

				if more {
					// set load offset for next query (pagination)
					after = fetchedIDs[len(fetchedIDs)-1].ID3()
				}

				for _, id := range fetchedIDs {
					// get full transaction by id
					tx, err := app.tonLiteClient.GetTransaction(context, shard, address.NewAddress(0, 0, id.Account), id.LT)
					if err != nil {
						app.logger.Error(errors.New("get tx:"+err.Error()), nil)
						return err
					}

					// add transaction to list
					txList = append(txList, tx)
				}
			}
		}

		for i, transaction := range txList {

			// check if transaction is related to any collection
			for _, collection := range collections {
				if hasCollection(collection.FriendlyAddress, txList) {
					payload, err := json.Marshal(collection.FriendlyAddress)
					if err != nil {
						app.logger.Error(fmt.Errorf("error: %s", err), nil)
						return err
					}
					runGetCollection := asynq.NewTask(database.TYPE_MIGRATE_COLLECTION, payload)

					info, err := app.asynqClient.Enqueue(runGetCollection, asynq.TaskID(collection.FriendlyAddress), asynq.MaxRetry(5), asynq.ProcessIn(5*time.Second), asynq.Retention(60 * time.Second), asynq.Queue(database.PRIORITY_URGENT))
					if err != nil {
						app.logger.Error(fmt.Errorf("error: %s", err), nil)
						return err
					}

					app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))
				}
			}

			app.logger.Info("processing transaction", "id", transaction.String(), "index", i)
		}

		if len(txList) == 0 {
			app.logger.Info("no transactions found")
		}

		master, err = app.tonLiteClient.WaitForBlock(master.SeqNo + 1).GetMasterchainInfo(context)
		if err != nil {
			app.logger.Error(errors.New("get masterchain info:"+err.Error()), nil)
			return err
		}
	}

}

func hasCollection(collectioAddr string, txs []*tlb.Transaction) bool {
	for _, tx := range txs {
		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal && (tx.IO.In.AsInternal().DstAddr.String() == collectioAddr || tx.IO.In.AsInternal().SrcAddr.String() == collectioAddr) {
			return true
		}
	}
	return false
}

func (app *application) runGoroutine(stop chan struct{}) bool {
	// Create a ticker to run health checks.
	ticker := time.NewTicker(time.Minute)

	var heartbeat time.Time

	// This is the main loop for the goroutine.
	for {
		select {
		// If we receive a signal to stop, we exit the goroutine.
		case <-stop:
			return true
		// If the ticker fires, we run a health check.
		case <-ticker.C:
			// Check the health of the goroutine. This could involve checking a global
			// "heartbeat" variable or other health checks.
			if !checkHealth(heartbeat) {
				// If the health check fails, log the error and consider exiting or restarting the goroutine.
				app.logger.Error(fmt.Errorf("goroutine is unhealthy"), nil)
				return false
			}
		// If neither stop nor ticker have triggered, we run the task.
		default:
			// Run the task and immediately update the heartbeat.
			// If the task panics or hangs, the defer recover() will catch it, and
			// the heartbeat will reflect the last successful run, not the failed one.
			err := app.RunListeningTransactions(context.Background())
			heartbeat = time.Now()

			if err != nil {
				// Log the error and consider restarting the goroutine.
				app.logger.Error(err, nil)
				return false
			}

			// After running the task, sleep for a while to avoid hammering the API.
			time.Sleep(time.Second * 10)
		}
	}
}
