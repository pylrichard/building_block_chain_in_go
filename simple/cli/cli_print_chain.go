package cli

import (
	"fmt"
	"strconv"

	"github.com/pylrichard/building_block_chain_in_go/simple/block"
)

func (cli *CLI) printChain(nodeId string) {
	bc := block.NewChain(nodeId)
	defer bc.Db.Close()

	bci := bc.Iterator()

	for {
		b := bci.Next()
		fmt.Printf("---- Block %x", b.Hash)
		fmt.Printf("Height: %d\n", b.Height)
		fmt.Printf("Prev Block: %x\n", b.PrevBlockHash)
		pow := block.NewProofOfWork(b)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range b.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
}