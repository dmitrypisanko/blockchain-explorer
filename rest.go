package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"strconv"
	"github.com/davecgh/go-spew/spew"
	//sq "gopkg.in/Masterminds/squirrel.v1"
	//"strings"
	//"fmt"
	//"strings"
	//"math/big"
	"log"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
	"github.com/mailru/dbr"
)

func (rest *Rest) run() {
	//var err error

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	/*
	r.GET("/blocks/", func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "20")
		offset := c.DefaultQuery("offset", "0")
		sortField := c.DefaultQuery("sortField", "number")
		sortDirection := c.DefaultQuery("sortDirection", "desc")

		request := sq.Select("*").From("blocks").OrderBy(strings.Join([]string{sortField, sortDirection}, " "))
		sql, _, err := request.ToSql()

		blocks := []TBlock{}
		err = rest.db.Select(&blocks, strings.Join([]string{sql, "limit", offset, ",", limit}, " "))

		spew.Dump(sql)
		spew.Dump(strings.Join([]string{sql, "limit", offset, ",", limit}, " "))
		spew.Dump(err)

		c.JSON(http.StatusOK, gin.H{
			"data": blocks,
		})
	})

	r.GET("/block/:number", func(c *gin.Context) {
		number, _ := strconv.Atoi(c.Param("number"))

		block := TBlock{}
		err = rest.db.Get(&block, "SELECT * FROM blocks where number = ?", number)

		c.JSON(http.StatusOK, gin.H{
			"data": block,
		})
	})

	r.GET("/transaction/:hash", func(c *gin.Context) {
		hash := c.Param("hash")

		transaction := TTransaction{}
		err = rest.db.Get(&transaction, "SELECT * FROM transactions where hash = ?", hash)

		c.JSON(http.StatusOK, gin.H{
			"data": transaction,
		})
	})
	*/

	r.GET("/transactions/", func(c *gin.Context) {
		/*
		limit := c.DefaultQuery("limit", "20")
		offset := c.DefaultQuery("offset", "0")
		sortField := c.DefaultQuery("sortField", "timestamp")
		sortDirection := c.DefaultQuery("sortDirection", "desc")
		accountAddress := c.DefaultQuery("account", "")
		blockNumber := c.DefaultQuery("blockNumber", "")

		request := sq.Select("*").From("transactions").OrderBy(strings.Join([]string{sortField, sortDirection}, " ")).Where(sq.Eq{"blockNumber": "1000"})

		if accountAddress != "" {
			request = request.Where(sq.Or{sq.Eq{"from": accountAddress}, sq.Eq{"to": accountAddress}})
		}

		if blockNumber != "" {
			request = request.Where("blockNumber = ?", blockNumber)
		}

		sql, _, err := request.ToSql()

		spew.Dump(accountAddress)
		spew.Dump(blockNumber)
		spew.Dump(sql)
		spew.Dump(strings.Join([]string{sql, "limit", offset, ",", limit}, " "))
		spew.Dump(err)

		transactions := []TTransaction{}
		err = rest.db.Select(&transactions, strings.Join([]string{sql, "limit", offset, ",", limit}, " "))
		*/

		transactions := []TTransaction{}
		query := rest.db.Select("*").From("transactions")
		query.Where(dbr.Eq("blockNumber", 457872))
		if _, err := query.Load(&transactions); err != nil {
			log.Fatal(err)
		}

		spew.Dump(transactions)

		//for _, item := range transactions {
		//	//log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_time: %s", item.CountryCode, item.OsID, item.BrowserID, item.Categories, item.ActionTime)
		//}

		c.JSON(http.StatusOK, gin.H{
			"data": transactions,
		})
	})

	r.GET("/account/:address", func(c *gin.Context) {
		address := c.Param("address")

		balance, err := rest.etherClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			log.Fatal(err)
		}

		nonce, err := rest.etherClient.NonceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			log.Fatal(err)
		}

		bytecode, err := rest.etherClient.CodeAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			log.Fatal(err)
		}

		isContract := false
		if len(bytecode) > 0 {
			isContract = true
		}

		c.JSON(http.StatusOK, gin.H{
			"data": &TAccount{
				Address: address,
				Balance: balance.String(),
				TransactionsCount: nonce,
				IsContract: isContract,
				Code: hex.EncodeToString(bytecode),
			},
		})

		/*
		transactions := []TTransaction{}
		err = rest.db.Select(&transactions, "SELECT * FROM transactions where from = ?", address)

		c.JSON(http.StatusOK, gin.H{
			"data": transactions,
		})
		*/
	})

	go r.Run()
}

