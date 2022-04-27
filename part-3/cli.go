package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const usage = `
Usage:
	add_block -data BLOCK_DATA	add a block to the block_chain
	print_chain					print all the blocks of the block_chain
`

func (cli *CLI) printUsage() {
	fmt.Printf(usage)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add_block", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print_chain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")

	switch os.Args[1] {
	case "add_block":
		err := addBlockCmd.Parse(os.Args[2:])
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

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.bc.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}