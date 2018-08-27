package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/jmoiron/sqlx"
	"github.com/mailru/dbr"
	"go.uber.org/zap"
)

type TConfig struct {
	NodeUrl string `yaml:"node_url"`
	Clickhouse string `yaml:"clickhouse"`
	Threads int `yaml:"threads"`
	BlockInterval int `yaml:"blockInterval"`
}

type TTransaction struct {
	Date string 		`db:"date"`
	Hash string 		`db:"hash"`
	BlockNumber uint64 	`db:"blockNumber"`
	Value string 		`db:"value"`
	GasUsed uint64 		`db:"gasUsed"`
	GasPrice uint64 	`db:"gasPrice"`
	Nonce uint64 		`db:"nonce"`
	To string 			`db:"to"`
	From string 		`db:"from"`
	Timestamp uint64 	`db:"timestamp"`
	//Data []byte
}

type TBlock struct {
	Date string					`db:"date"`
	Number uint64				`db:"number"`
	Timestamp uint64			`db:"timestamp"`
	Hash string					`db:"hash"`
	ParentHash string			`db:"parentHash"`
	UncleHash string			`db:"uncleHash"`
	MinedBy string				`db:"minedBy"`
	GasUsed uint64				`db:"gasUsed"`
	GasLimit uint64				`db:"gasLimit"`
	Nonce uint64				`db:"nonce"`
	Size float64				`db:"size"`
	TransactionsCount uint64	`db:"transactionsCount"`
	Difficulty uint64			`db:"difficulty"`
	Extra string				`db:"extra"`
}

type TAccount struct {
	Address string
	Balance string
	TransactionsCount uint64
	IsContract bool
	Code string
}

type Watcher struct {
	etherClient *ethclient.Client
	//db *sqlx.DB
	db *dbr.Session
	logger *zap.Logger
	lastParsedBlockNumber int64
	lastBlockNumber int64
	queue chan int64
}

type Rest struct {
	//db *sqlx.DB
	db *dbr.Session
	etherClient *ethclient.Client
}