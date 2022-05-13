package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/pylrichard/building_block_chain_in_go/simple/block"
	"github.com/pylrichard/building_block_chain_in_go/simple/transaction"
)

const protocol = "tcp"
const nodeVersion = 1
const cmdLen = 12

var nodeAddr string
var miningAddr string
var knownNodes = []string{"localHost:3000"}
var blocksInTransit [][]byte
var memPool = make(map[string]transaction.Transaction)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom	string
	Block		[]byte
}

type GetBlocks struct {
	AddrFrom string
}

type GetData struct {
	AddrFrom	string
	Type		string
	Id			[]byte
}

type Inv struct {
	AddrFrom	string
	Type		string
	Items		[][]byte
}

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

func bytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
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

func sendGetBlocks(addr string) {
	payload := gobEncode(GetBlocks{nodeAddr})
	request := append(cmdToBytes("get_blocks"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *block.Chain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{nodeVersion, bestHeight, nodeAddr})
	request := append(cmdToBytes("version"), payload...)

	sendData(addr, request)
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[cmdLen:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))

	requestBlocks()
}

func handleConnection(conn net.Conn, bc *block.Chain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	cmd := bytesToCmd(request[:cmdLen])
	fmt.Printf("Recieved %s command\n", cmd)

	switch cmd {
	case "addr":
		handleAddr(request)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
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

func isNodeKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}