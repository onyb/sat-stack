package main

import (
	"log"
	"os"

	"ledger-sats-stack/pkg/handlers"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/gin-gonic/gin"
)

func main() {
	client := getRPCClient()
	engine := getRouter(client)
	defer client.Shutdown()
	engine.Run()
}

func getRouter(client *rpcclient.Client) *gin.Engine {
	engine := gin.Default()

	engine.GET("/blockchain/v3/blocks/:block", handlers.GetBlock(client))
	engine.GET("/blockchain/v3/transactions/:hash", handlers.GetTransaction(client))
	engine.GET("/blockchain/v3/transactions/:hash/hex", handlers.GetTransactionHex(client))

	engine.GET("/blockchain/v3/explorer/_health", handlers.GetHealth(client))

	return engine
}

func getRPCClient() *rpcclient.Client {
	connCfg := &rpcclient.ConnConfig{
		Host:         os.Getenv("BITCOIND_RPC_HOST"),
		User:         os.Getenv("BITCOIND_RPC_USER"),
		Pass:         os.Getenv("BITCOIND_RPC_PASSWORD"),
		HTTPPostMode: true,
		DisableTLS:   os.Getenv("BITCOIND_RPC_ENABLE_TLS") != "true",
	}
	// The notification parameter is nil since notifications are not
	// supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}