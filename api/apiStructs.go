package api

import (
	"github.com/zachvanuum/tarkus/blockchain"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ChainResponse struct {
	Chain []blockchain.Block `json:"chain"`
	Length int `json:"length"`
}

type MineResponse struct {
	Message string `json:"message"`
	Index int64 `json:"index"`
	Transactions []blockchain.Transaction `json:"transactions"`
	Proof int64 `json:"proof"`
	PreviousHash string `json:"previousHash"`
}

type RegisterNodeRequest struct {
	URL string `json:"string"`
}
