package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
	"github.com/mailru/dbr"
	"strconv"
	"strings"
)

func (rest *Rest) run() {
	//var err error

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/blocks/", func(c *gin.Context) {
		limit, _ := strconv.ParseUint(c.DefaultQuery("limit", "20"), 10, 64)
		offset, _ := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
		sortField := c.DefaultQuery("sortField", "timestamp")
		sortDirection := c.DefaultQuery("sortDirection", "desc")

		blocks := []TBlock{}
		query := rest.db.Select("*").From("blocks").Limit(limit).Offset(offset).OrderBy(strings.Join([]string{sortField, sortDirection}, " "))

		if _, err := query.Load(&blocks); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": blocks,
		})
	})

	r.GET("/block/:number", func(c *gin.Context) {
		number, _ := strconv.Atoi(c.Param("number"))

		block := TBlock{}
		query := rest.db.Select("*").From("blocks").Where("number = ?", number)
		if _, err := query.Load(&block); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": block,
		})
	})

	r.GET("/transaction/:hash", func(c *gin.Context) {
		hash := c.Param("hash")

		transaction := TTransaction{}
		query := rest.db.Select("*").From("transactions").Where("hash = ?", hash)
		if _, err := query.Load(&transaction); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": transaction,
		})
	})


	r.GET("/transactions/", func(c *gin.Context) {
		limit, _ := strconv.ParseUint(c.DefaultQuery("limit", "20"), 10, 64)
		offset, _ := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
		sortField := c.DefaultQuery("sortField", "timestamp")
		blockNumber, _ := strconv.ParseUint(c.DefaultQuery("blockNumber", ""), 10, 64)
		accountAddress := c.DefaultQuery("account", "")
		sortDirection := c.DefaultQuery("sortDirection", "desc")

		transactions := []TTransaction{}
		query := rest.db.Select("*").From("transactions").Limit(limit).Offset(offset).OrderBy(strings.Join([]string{sortField, sortDirection}, " "))

		if accountAddress != "" {
			query = query.Where(dbr.Or(dbr.Eq("from", accountAddress), dbr.Eq("to", accountAddress)))
		}

		if blockNumber != 0 {
			query = query.Where(dbr.Eq("blockNumber", blockNumber))
		}

		if _, err := query.Load(&transactions); err != nil {
			log.Fatal(err)
		}

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
	})

	go r.Run()
}

