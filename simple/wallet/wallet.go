package wallet

import (
	"bytes"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"

	"github.com/pylrichard/building_block_chain_in_go/simple/codec"
)

const addrChecksumLen = 4

func HashPubKey(pubKey []byte) []byte {
	pubSha256 := sha256.Sum256(pubKey)

	ripemd160Hash := ripemd160.New()
	_, err := ripemd160Hash.Write(pubSha256[:])
	if err != nil {
		log.Panic(err)
	}
	pubRipemd160 := ripemd160Hash.Sum(nil)

	return pubRipemd160
}

func ValidateAddr(addr string) bool {
	pubKeyHash := codec.Base58Decode([]byte(addr))
	actualChecksum := pubKeyHash[len(pubKeyHash) - addrChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - addrChecksumLen]
	targetChecksum := getChecksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func getChecksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addrChecksumLen]
}