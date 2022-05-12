package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pylrichard/building_block_chain_in_go/simple/block"
)

const protocol = "tcp"
const nodeVersion = 1
const cmdLen = 12

var nodeAddr string
var miningAddr string
var knownNodes = []string{"localHost:3000"}


type Tx struct {
	AddrFrom	string
	Transaction []byte
}

type Version struct {
	Version		int
	BestHeight	int
	AddrFrom	string
}

func StartServer(nodeId, minerAddr string) {
	nodeAddr = fmt.Sprintf("localHost: %s", nodeId)
	miningAddr = minerAddr

	l, err := net.Listen(protocol, nodeAddr)
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()

	bc := block.NewChain(nodeId)

	if nodeAddr != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func cmdToBytes(cmd string) []byte {
	var bytes [cmdLen]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updateNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updateNodes = append(updateNodes, node)
			}
		}

		knownNodes = updateNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendVersion(addr string, bc *block.Chain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{nodeVersion, bestHeight, nodeAddr})
	request := append(cmdToBytes("version"), payload...)

	sendData(addr, request)
}

func handleConnection(conn net.Conn, bc *block.Chain) {

}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}