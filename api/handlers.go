package api

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/zachvanuum/tarkus/blockchain"
)

// TODO: Use Logrus for logging request/response to http server

func GetChainHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[GetChainHandler] - Received request")

		response := ChainResponse{
			Chain: chain.Chain,
			Length: len(chain.Chain),
		}

		if err := unmarshalToResponse(w, response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[GetChainHandler] - Failed to unmarshal GetChain response", err)
		}
	}
}

func GetMineHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[GetMineHandler] - Received request")

		lastBlock := chain.LastBlock()
		proof := chain.ProofOfWork(lastBlock.Proof)
		prevHash := lastBlock.PreviousHash

		// Sender is 0 to signify that the this node mined a new coin
		chain.NewTransaction("0", chain.ID, 1)

		block := chain.NewBlock(proof, prevHash)
		response := MineResponse{
			Message: "A new block was mined.",
			Index: block.Index,
			Transactions: block.Transactions,
			Proof: block.Proof,
			PreviousHash: block.PreviousHash,
		}

		if err := unmarshalToResponse(w, response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[GetMineHandler] - Failed to unmarshal GetMine response", err)
		}		
	}
}

func GetConsensusHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newChain []blockchain.Block
		maxLen := len(chain.Chain)
	
		for node, _ := range chain.Nodes {
			res, err := http.Get(node + "/chain")
			if err != nil {
				log.Printf("[ResolveConflicts] - Failed to retrieve chain for node %s during consensus.\n", node)
			}
					
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				decoder := json.NewDecoder(res.Body)
				var chainRes ChainResponse
	
				if err := decoder.Decode(&chainRes); err == nil {
					if chainRes.Length > maxLen && blockchain.ValidChain(chain.Chain) {
						maxLen = chainRes.Length
						newChain = chainRes.Chain
					}
				} else {
					log.Printf("[ResolveConflicts] - Failed to decode chain response from node %s\n", node)
				}	
			}
	
		}
	
		if len(newChain) > 0 {
			chain.Chain = newChain
			// return true
		}
	
		// return false
	}
}

func PostNewTransactionHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[PostNewTransactionHandler] - Received request")

		decoder := json.NewDecoder(r.Body)
		var transaction blockchain.Transaction

		if err := decoder.Decode(&transaction); err == nil {
			index := chain.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)
			response := MessageResponse{ Message: fmt.Sprintf("Transaction added to block at index %d", index) }

			if err := unmarshalToResponse(w, response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[PostNewTransactionHandler] - Failed to unmarshal PostNewTransaction response", err)
			}
		} else {
			http.Error(w, "Request body missing transaction fields: expected a sender, recipient, and famount.", http.StatusBadRequest)
			log.Println("[PostNewTransactionHandler] - Did not receive needed fields in POST body, failed to decode body to Transaction", err)
		}
		defer r.Body.Close()
	}
}

func PostRegisterNodeHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[PostRegisterNodeHandler] - Received request")
		
		decoder := json.NewDecoder(r.Body)
		var reqBody RegisterNodeRequest
		
		if err := decoder.Decode(&reqBody); err == nil {
			chain.RegisterNode(reqBody.URL)

			w.WriteHeader(http.StatusOK)
			log.Printf("[PostRegiterNodeHandler] - Registered new node %s\n", reqBody.URL)
		} else {
			http.Error(w, "Request body missing url field.", http.StatusBadRequest)
			log.Println("[PostRegisterNodeHandler] - Did not receive needed fields in POST body, failed to decode body to RegisterNodeBody", err)
		}
	}
}

func unmarshalToResponse(w http.ResponseWriter, response interface{}) error {
	var jsonResponse []byte
	var err error

	if jsonResponse, err = json.Marshal(response); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}

	return err
}
