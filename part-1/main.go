package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("BTC 1")
	bc.AddBlock("BTC 2")

	for _, b := range bc.blocks {
		fmt.Printf("Prev hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Println()
	}
}