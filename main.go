package main

// export GOPATH=$HOME/gocode

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/davecgh/go-spew/spew"

	//_ "github.com/kshvakov/clickhouse"
	//"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	_ "github.com/mailru/go-clickhouse"
	"github.com/mailru/dbr"
)


func loadConf() TConfig {
	var config TConfig

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
}


func main() {
	config := loadConf()

	//db, err := sqlx.Connect("clickhouse", config.Clickhouse)
	connect, err := dbr.Open("clickhouse", config.Clickhouse, nil)
	if err != nil {
		log.Fatalln(err)
	}

	db := connect.NewSession(nil)

	client, err := ethclient.Dial(config.NodeUrl)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	defer logger.Sync()

	rest := &Rest{
		db: db,
		etherClient: client,
	}
	rest.run()

	watcher := &Watcher{
		lastBlockNumber: 0,
		lastParsedBlockNumber: 0,
		etherClient: client,
		logger: logger,
		db: db,
	}

	watcher.run(config.Threads)
	select {}
}