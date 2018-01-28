package blockchain

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	Blockchain *Blockchain
}

func (cli *CLI) Usage() {
	fmt.Println(`
Available Commands:
  add-block      -
  print-chain    -    Print all the blocks of the blockchain
`)
}

func (cli *CLI) AddBlock(data string) {
	if err := cli.Blockchain.AddBlock(data); err != nil {
		panic(err)
	}

	fmt.Println("Added new block")
}

func (cli *CLI) PrintChain() {
	bci := cli.Blockchain.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %v\n", pow.Validate())

		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	if err := cli.Validate(); err != nil {
		cli.Usage()
		os.Exit(1)
	}

	addBlockCmd := flag.NewFlagSet("add-block", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print-chain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	var err error

	switch os.Args[1] {
	case "add-block":
		err = addBlockCmd.Parse(os.Args[2:])
	case "print-chain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.Usage()
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		cli.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *CLI) Validate() error {
	if len(os.Args) < 2 {
		return errors.New("Command not specified")
	}

	return nil
}
