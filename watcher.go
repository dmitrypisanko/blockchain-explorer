package main

import (
	"log"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"time"
	"github.com/vjeantet/jodaTime"
	"math/big"
	"go.uber.org/zap"
)

func (watcher *Watcher) parseBlock(blockNumber int64) {
	//watcher.logger.Info("Parse block", zap.Int64("number", blockNumber))

	block, err := watcher.etherClient.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal(err)
	}

	timestamp := block.Time().Int64()

	// dirty hack
	if timestamp == 0 {
		timestamp = 1438269975
	}

	date := jodaTime.Format("YYYY-MM-dd", time.Unix(timestamp, 0))

	transactions := []TTransaction{}
	for _, tx := range block.Transactions() {
		if msg, err := tx.AsMessage(types.HomesteadSigner{}); err == nil {
			to := ""
			if tx.To() != nil {
				to = tx.To().Hex()
			}

			transactions = append(transactions, TTransaction{
				Date: 			date,
				BlockNumber: 	block.NumberU64(),
				Hash:        	tx.Hash().Hex(),
				Value:       	tx.Value().String(),
				GasUsed:     	tx.Gas(),
				GasPrice:    	tx.GasPrice().Uint64(),
				Nonce:       	tx.Nonce(),
				To:        		to,
				From:      		msg.From().Hex(),
				Timestamp: 		block.Time().Uint64(),
				//Data: tx.Data(),
			})

			//dbTX := watcher.db.MustBegin()
			//dbTX.MustExec(
			//	"INSERT INTO transactions (date, timestamp, hash, blockNumber, value, gasUsed, gasPrice, nonce, to, from) VALUES ($1, $2, $3, $4, $5, $6, $6, $7, $7, $8, $9, $10)",
			//	date, block.Time().Uint64(), tx.Hash().Hex(), block.NumberU64(), tx.Value().String(), tx.Gas(), tx.GasPrice().Uint64(), tx.Nonce(), to, msg.From().Hex(),
			//)
			//dbTX.Commit()
		}
	}

	//dbTX := watcher.db.MustBegin()
	//dbTX.MustExec(
	//	"INSERT INTO blocks (date, timestamp, hash, number, gasUsed, gasLimit, nonce, size, transactionsCount, difficulty, extra, parentHash, uncleHash, minedBy) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
	//	date, timestamp, block.Hash().Hex(), block.NumberU64(), block.GasUsed(), block.GasLimit(), block.Nonce(), block.Size(), len(block.Transactions()), block.Difficulty().Uint64(), string(block.Extra()[:]), block.ParentHash().Hex(), block.UncleHash().Hex(), block.Coinbase().Hex(),
	//)
	//dbTX.Commit()
}

func (watcher *Watcher) queueWatcher() {
	for {
		blockNumber := <-watcher.queue
		watcher.parseBlock(blockNumber)
	}
}

func (watcher *Watcher) getLastBlock() {
	header, err := watcher.etherClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	watcher.lastBlockNumber = header.Number.Int64()
	watcher.logger.Info("getLastBlock", zap.Int64("lastBlockNumber", watcher.lastBlockNumber))
}


func (watcher *Watcher) blockWatcher(timeout int) {
	for {
		watcher.getLastBlock()

		for i := watcher.lastParsedBlockNumber + 1; i <= watcher.lastBlockNumber; i++ {
			watcher.logger.Info("newBlock", zap.Int64("number", i))
			watcher.queue <- i
		}

		watcher.lastParsedBlockNumber = watcher.lastBlockNumber
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

func (watcher *Watcher) run(threads int) {
	watcher.getLastBlock()

	//var blocks []struct {
	//	Number int64 `db:"number"`
	//}

	////if err := watcher.db.Select(&blocks, "SELECT number FROM blocks order by number asc"); err != nil {
	////	log.Fatal(err)
	////}
	//
	//watcher.logger.Info("Blocks in DB", zap.Int("number", len(blocks)))
	//
	blockHash := map[int64]int {}
	//for _, item := range blocks {
	//	blockHash[item.Number] = 1
	//}

	watcher.queue = make(chan int64, threads)
	for w := 0; w < threads; w++ {
		go watcher.queueWatcher()
	}

	for i := int64(0); i < watcher.lastBlockNumber; i++ {
		_, ok := blockHash[i]
		if !ok {
			if i%10000 == 0 {
				//watcher.logger.Info("Send block to queue", zap.Int64("number", i))
			}

			//api.queue <- i
		}

		watcher.lastParsedBlockNumber = int64(i)
	}

	//go watcher.blockWatcher(config.BlockInterval)
}
