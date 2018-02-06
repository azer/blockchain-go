package blockchain

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CLI struct{}

func (cli *CLI) Usage() {
	fmt.Println(`
Available Commands:
  get-balance          -    Get balance of <address>
  create-blockchain    -    Create a blockchain and send genesis reward to <address>
  send                 -    Send N coins from one to another. (e.g send -from <x> -to <y> -amont <z>)
  print-chain          -    Print all the blocks of the blockchain
`)
}

func (cli *CLI) CreateBlockchain() {
	cmd := flag.NewFlagSet("create-blockchain", flag.ExitOnError)
	rewardAddress := cmd.String("address", "", "The address to send genesis block reward to.")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	if *rewardAddress == "" {
		panic("Reward address required")
	}

	blockchain, err := NewBlockchain(*rewardAddress)
	if err != nil {
		panic(err)
	}

	blockchain.DB.Close()
	fmt.Println("Done")
}

func (cli *CLI) GetBalance() {
	cmd := flag.NewFlagSet("get-balance", flag.ExitOnError)
	address := cmd.String("address", "", "The address to get balance for")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	if *address == "" {
		panic("Wallet address required")
	}

	bc, err := RestoreBlockchain()
	if err != nil {
		panic(err)
	}

	defer bc.DB.Close()

	balance := 0
	UTXOs, err := bc.FindUTXO(*address)
	if err != nil {
		panic(err)
	}

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Println(fmt.Sprintf("Balance of '%s': %d", *address, balance))
}

func (cli *CLI) Send() {
	cmd := flag.NewFlagSet("get-balance", flag.ExitOnError)
	from := cmd.String("from", "", "Source wallet address")
	to := cmd.String("to", "", "Destination wallet address")
	amount := cmd.Int("amount", 0, "Amount to send")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	bc, err := RestoreBlockchain()
	if err != nil {
		panic(err)
	}

	defer bc.DB.Close()

	tx, err := NewUTXOTransaction(*from, *to, *amount, bc)
	if err != nil {
		panic(err)
	}

	bc.MineBlock([]*Transaction{tx})

	fmt.Println("Success!")
}

func (cli *CLI) PrintBlockchain() {
	bc, err := RestoreBlockchain()
	if err != nil {
		panic(err)
	}

	defer bc.DB.Close()
	bci := bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
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

	switch os.Args[1] {
	case "get-balance":
		cli.GetBalance()
	case "create-blockchain":
		cli.CreateBlockchain()
	case "print-chain":
		cli.PrintBlockchain()
	case "send":
		cli.Send()
	default:
		cli.Usage()
		os.Exit(1)
	}
}

func (cli *CLI) Validate() error {
	if len(os.Args) < 2 {
		return errors.New("Command not specified")
	}

	return nil
}
