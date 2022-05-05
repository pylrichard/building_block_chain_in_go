package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const argsNum = 2

type CLI struct {}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  create_block_chain -address ADDRESS - Create a block_chain and send genesis block reward to ADDRESS")
	fmt.Println("  create_wallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  get_balance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  list_addresses - Lists all addresses from the wallet file")
	fmt.Println("  print_chain - Print all the blocks of the block_chain")
	fmt.Println("  reindex_utxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  start_node -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < argsNum {
		cli.printUsage()
		os.Exit(1)
	}
}

//Run 解析命令参数并执行命令
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeId := os.Getenv("NODE_ID")
	if nodeId == "" {
		fmt.Printf("NODE_ID is not set!")
		os.Exit(1)
	}

	getBalanceCmd := flag.NewFlagSet("get_balance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("create_block_chain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("create_wallet", flag.ExitOnError)
	listAddrCmd := flag.NewFlagSet("list_addr", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print_chain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindex_utxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("start_node", flag.ExitOnError)

	getBalanceAddr := getBalanceCmd.String("addr", "", "The address to get balance for")
	createBlockChainAddr := createBlockChainCmd.String("addr", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")

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
	case "create_wallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "list_addr":
		err := listAddrCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print_chain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindex_utxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "start_node":
		err := startNodeCmd.Parse(os.Args[2:])
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
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddr == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
	}

	if createWalletCmd.Parsed() {
	}

	if listAddrCmd.Parsed() {
	}

	if printChainCmd.Parsed() {
	}

	if reindexUTXOCmd.Parsed() {
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount, nodeId, *sendMine)
	}

	if startNodeCmd.Parsed() {
		cli.startNode(nodeId, *startNodeMiner)
	}
}