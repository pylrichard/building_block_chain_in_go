package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {}

func (cli *CLI) createBlockChain(addr string) {
	bc := CreateBlockChain(addr)
	bc.db.Close()
	fmt.Println("Create block chain done!")
}

func (cli *CLI) getBalance(addr string) {
	bc := NewBlockChain()
	defer bc.db.Close()

	balance := 0
	utxs := bc.FindUTX(addr)

	for _, out := range utxs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", addr, balance)
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("    get_balance -addr ADDRESS - Get balance of ADDRESS")
	fmt.Println("    create_block_chain -addr ADDRESS - Create a block chain and send genesis block reward to ADDRESS")
	fmt.Println("    print_chain - Print all the blocks of the block chain")
	fmt.Println("    send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bc := NewBlockChain()
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		b := bci.Next()

		fmt.Printf("Prev hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Hash: %x\n", b.Hash)
		pow := NewProofOfWork(b)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Send successfully!")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("get_balance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("create_block_chain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print_chain", flag.ExitOnError)

	getBalanceAddr := getBalanceCmd.String("addr", "", "The address to get balance for")
	createBlockChainAddr := createBlockChainCmd.String("addr", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "get_balance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "create_block_chain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print_chain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddr == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddr)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddr == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createBlockChainAddr)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}