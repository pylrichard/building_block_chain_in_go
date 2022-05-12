package cli

import (
	"fmt"
	"log"

	"github.com/pylrichard/building_block_chain_in_go/simple/server"
	"github.com/pylrichard/building_block_chain_in_go/simple/wallet"
)

func (cli *CLI) startNode(nodeId, minerAddr string) {
	fmt.Printf("Starting node %s\n", nodeId)
	if len(minerAddr) > 0 {
		if wallet.ValidateAddr(minerAddr) {
			fmt.Println("Mining is on, Address to receive rewards:", minerAddr)
		} else {
			log.Panic("Wrong miner Address")
		}
	}

	server.StartServer(nodeId, minerAddr)
}