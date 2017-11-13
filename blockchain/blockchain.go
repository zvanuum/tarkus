package blockchain

import (
	"fmt"
	"time"
	"crypto/sha256"

	"github.com/satori/go.uuid"
)

func InitializeBlockchain() Blockchain {
	blockchain := Blockchain{
		Chain: make([]Block, 0),
        CurrentTransactions: make([]Transaction, 0),
	}

	blockchain.ID = uuid.NewV4().String()
	blockchain.Nodes = make(map[string]bool)
	blockchain.NewBlock(100, "1")
	return blockchain
}


type BlockchainApp interface {
	NewBlock(proof int, previousHash string) Block
	NewTransaction(sender string, recipient string, amount float64) Block
	LastBlock() Block
	ProofOfWork(lastProof int64) int64 
	RegisterNode(url string)
	ResolveConflicts() bool
}

type Blockchain struct {
	ID string `json:"id"`
	Nodes map[string]bool `json:"nodes"`
	Chain []Block `json:"chain"`
	CurrentTransactions []Transaction `json:"currentTransactions"`
}

func (blockchain *Blockchain) NewBlock(proof int64, previousHash string) Block {
	block := Block{
		Index: int64(len(blockchain.Chain) + 1),
		Timestamp: time.Now().Unix(),
		Transactions: blockchain.CurrentTransactions,
		Proof: proof,
		PreviousHash: previousHash,
	}

	blockchain.CurrentTransactions = make([]Transaction, 0)
	blockchain.Chain = append(blockchain.Chain, block)

	return block
}

func (blockchain *Blockchain) NewTransaction(sender string, recipient string, amount float64) int64 {
	transaction := Transaction{ Sender: sender, Recipient: recipient, Amount: amount}
	blockchain.CurrentTransactions = append(blockchain.CurrentTransactions, transaction)

	return blockchain.LastBlock().Index + 1
}

func (blockchain *Blockchain) LastBlock() Block {
	return blockchain.Chain[len(blockchain.Chain) - 1]
}

func (blockchain *Blockchain) ProofOfWork(lastProof int64) int64 {
	proof := int64(0)

	for !ValidProof(lastProof, proof) {
		proof += 1
	}

	return proof
}

func (blockchain *Blockchain) RegisterNode(url string) {
	blockchain.Nodes[url] = true
}

func ValidProof(lastProof int64, proof int64) bool {
	proofStr := fmt.Sprintf("%d%d", lastProof, proof)
	hashed := sha256.Sum256([]byte(proofStr))
	return fmt.Sprintf("%x", hashed)[:4] == "0000"
}

func  ValidChain(blockchain []Block) bool {
	prevBlock := blockchain[0]
	
	for i := 1; i < len(blockchain); i++ {
		block := blockchain[i]

		if block.PreviousHash != prevBlock.Hash() {
			return false
		}

		if !ValidProof(prevBlock.Proof, block.Proof) {
			return false
		}

		prevBlock = block
	}

	return true
}
