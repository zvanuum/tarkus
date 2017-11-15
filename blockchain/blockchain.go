package blockchain

import (
	"log"
	"fmt"
	"time"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/zachvanuum/tarkus/models"
)

func InitializeBlockchain() Blockchain {
	blockchain := Blockchain{
		Chain: make([]models.Block, 0),
        CurrentTransactions: make([]models.Transaction, 0),
	}

	blockchain.ID = uuid.NewV4().String()
	blockchain.Nodes = make(map[string]bool)
	blockchain.NewBlock(100, "1")
	return blockchain
}


type BlockchainApp interface {
	NewBlock(proof int, previousHash string) models.Block
	NewTransaction(sender string, recipient string, amount float64) models.Block
	LastBlock() models.Block
	ProofOfWork(lastProof int64) int64 
	RegisterNode(url string)
	ResolveConflicts() bool
	GetNodes() []string
}

type Blockchain struct {
	ID string `json:"id"`
	Nodes map[string]bool `json:"nodes"`
	Chain []models.Block `json:"chain"`
	CurrentTransactions []models.Transaction `json:"currentTransactions"`
}

func (blockchain *Blockchain) NewBlock(proof int64, previousHash string) models.Block {
	block := models.Block{
		Index: int64(len(blockchain.Chain) + 1),
		Timestamp: time.Now().Unix(),
		Transactions: blockchain.CurrentTransactions,
		Proof: proof,
		PreviousHash: previousHash,
	}

	blockchain.CurrentTransactions = make([]models.Transaction, 0)
	blockchain.Chain = append(blockchain.Chain, block)

	return block
}

func (blockchain *Blockchain) NewTransaction(sender string, recipient string, amount float64) int64 {
	transaction := models.Transaction{ Sender: sender, Recipient: recipient, Amount: amount}
	blockchain.CurrentTransactions = append(blockchain.CurrentTransactions, transaction)

	return blockchain.LastBlock().Index + 1
}

func (blockchain *Blockchain) LastBlock() models.Block {
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

func (blockchain *Blockchain) ResolveConflicts() bool {
	var newChain []models.Block
	maxLen := len(blockchain.Chain)

	for node, _ := range blockchain.Nodes {
		res, err := http.Get(node + "/chain")
		if err != nil {
			log.Printf("[ResolveConflicts] - Failed to retrieve chain for node %s during consensus.\n", node)
		}
				
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			decoder := json.NewDecoder(res.Body)
			var chainRes models.ChainResponse

			if err := decoder.Decode(&chainRes); err == nil {
				if chainRes.Length > maxLen && ValidChain(blockchain.Chain) {
					maxLen = chainRes.Length
					newChain = chainRes.Chain
				}
			} else {
				log.Printf("[ResolveConflicts] - Failed to decode chain response from node %s\n", node)
			}	
		}

	}

	if len(newChain) > 0 {
		blockchain.Chain = newChain
		return true
	}

	return false
}

func (blockchain *Blockchain) GetNodes() []string {
	nodes := make([]string, len(blockchain.Nodes))

	i := 0
	for node := range blockchain.Nodes {
		nodes[i] = node
		i++
	}

	return nodes
}

func ValidProof(lastProof int64, proof int64) bool {
	proofStr := fmt.Sprintf("%d%d", lastProof, proof)
	hashed := sha256.Sum256([]byte(proofStr))
	return fmt.Sprintf("%x", hashed)[:4] == "0000"
}

func  ValidChain(blockchain []models.Block) bool {
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
