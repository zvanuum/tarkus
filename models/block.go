package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

type Block struct {
	Index int64 `json:"index"`
	Timestamp int64 `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof int64 `json:"proof"`
	PreviousHash string `json:"previousHash"`
}

func (block *Block) Hash() string {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		log.Printf("Failed to hash block with index %d\n", block.Index)
		return ""
	}

	return fmt.Sprintf("%x", sha256.Sum256(blockJSON))
}
