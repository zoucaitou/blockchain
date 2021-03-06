package main

import (
  "bytes"
  "math"
	"math/big"
  "crypto/sha256"
)

var (
  maxNonce = math.MaxInt64
)

const targetBits = 1

type ProofOfWork struct {
  block *Block
  target *big.Int
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
  data := bytes.Join(
      [][]byte {
        pow.block.PrevBlockHash,
        pow.block.Data,
        IntToHex(pow.block.Timestamp),
        IntToHex(targetBits),
        IntToHex(int64(nonce)),
      },
      []byte{},
  )
  return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
  var hashInt big.Int
  var hash [32]byte

  nonce := 0

  logf("Mining the block containing %s", pow.block.Data)
  for nonce < maxNonce {
    data := pow.prepareData(nonce)
    hash = sha256.Sum256(data)

    logf("hash %s", hash)
    hashInt.SetBytes(hash[:])
    if hashInt.Cmp(pow.target) == -1 {
	    break
    } else {
	    nonce++
		}
  }
  return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
  var hashInt big.Int

  data := pow.prepareData(pow.block.Nonce)
  hash := sha256.Sum256(data)
  hashInt.SetBytes(hash[:])

  isValid := hashInt.Cmp(pow.target) == -1

  return isValid
}

func NewProofOfWork(b *Block) *ProofOfWork {
  target := big.NewInt(1)
  target.Lsh(target, uint(256 - targetBits))

  pow := &ProofOfWork{b, target}
  return pow
}
