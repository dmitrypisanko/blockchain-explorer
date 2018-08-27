package main

import (
	"log"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"time"
	"github.com/vjeantet/jodaTime"
	"math/big"
	"go.uber.org/zap"
	"github.com/davecgh/go-spew/spew"
)

func (watcher *Watcher) parseBlock(blockNumber int64) {
	watcher.logger.Info("Parse block", zap.Int64("number", blockNumber))

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

	tx, err := watcher.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.RollbackUnlessCommitted()

	for _, tx := range block.Transactions() {
		if msg, err := tx.AsMessage(types.HomesteadSigner{}); err == nil {
			to := ""
			if tx.To() != nil {
				to = tx.To().Hex()
			}

			transaction := TTransaction{
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
			}
			_, err = watcher.db.InsertInto("transactions").Columns("date", "timestamp", "hash", "blockNumber", "value", "gasUsed", "gasPrice", "nonce", "to", "from").Record(transaction).Exec()
			if err != nil {
				watcher.logger.Error("Insert transaction")
				spew.Dump(err)
			}
		}
	}

	blockInfo := TBlock{
		Date: date,
		Timestamp: timestamp,
		Hash: block.Hash().Hex(),
		Number: block.NumberU64(),
		GasUsed: block.GasUsed(),
		GasLimit: block.GasLimit(),
		Nonce :block.Nonce(),
		Size: block.Size().String(),
		TransactionsCount: len(block.Transactions()),
		Difficulty: block.Difficulty().Uint64(),
		Extra: string(block.Extra()[:]),
		ParentHash: block.ParentHash().Hex(),
		UncleHash: block.UncleHash().Hex(),
		MinedBy: block.Coinbase().Hex(),
	}

	_, err = watcher.db.InsertInto("blocks").Columns("date", "timestamp", "hash", "number", "gasUsed", "gasLimit", "nonce", "size", "transactionsCount", "difficulty", "extra", "parentHash", "uncleHash", "minedBy").Record(blockInfo).Exec()
	if err != nil {
		watcher.logger.Error("Insert block")
		spew.Dump(err)
	}

	tx.Commit()
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

func (watcher *Watcher) run(threads int, blockInterval int) {
	watcher.getLastBlock()

	var blocks []struct {
		Number int64 `db:"number"`
	}

	query := watcher.db.Select("*").From("blocks").OrderBy("number asc")
	if _, err := query.Load(&blocks); err != nil {
		log.Fatal(err)
	}

	watcher.logger.Info("Blocks in DB", zap.Int("number", len(blocks)))

	blockHash := map[int64]int {}
	for _, item := range blocks {
		blockHash[item.Number] = 1
	}

	watcher.queue = make(chan int64, threads)
	for w := 0; w < threads; w++ {
		go watcher.queueWatcher()
	}

	for i := int64(0); i < watcher.lastBlockNumber; i++ {
		_, ok := blockHash[i]
		if !ok {
			if i%10000 == 0 {
				watcher.logger.Info("Send block to queue", zap.Int64("number", i))
			}

			watcher.queue <- i
		}

		watcher.lastParsedBlockNumber = int64(i)
	}

	go watcher.blockWatcher(blockInterval)
}
