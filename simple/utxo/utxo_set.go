package utxo

import "github.com/pylrichard/building_block_chain_in_go/simple/block"

type Set struct {
	Chain *block.Chain
}

//Reindex 重新构建UTXO Set
func (u Set) Reindex() {

}