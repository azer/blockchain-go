package main

import (
	"github.com/azer/blockchain"
)

func main() {
	bc, err := blockchain.NewBlockChain()
	if err != nil {
		panic(err)
	}

	defer bc.DB.Close()

	cli := &blockchain.CLI{
		Blockchain: bc,
	}

	cli.Run()
}
