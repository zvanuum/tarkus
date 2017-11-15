package api

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/zachvanuum/tarkus/blockchain"
	"github.com/zachvanuum/tarkus/models"
)

func ChainHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ChainHandler] - Received request")

		response := models.ChainResponse{
			Chain: chain.Chain,
			Length: len(chain.Chain),
		}

		if err := unmarshalToResponse(w, response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ChainHandler] - Failed to unmarshal GetChain response", err)
		}
	}
}

func MineHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[MineHandler] - Received request")

		lastBlock := chain.LastBlock()
		proof := chain.ProofOfWork(lastBlock.Proof)
		prevHash := lastBlock.PreviousHash

		// Sender is 0 to signify that the this node mined a new coin
		chain.NewTransaction("0", chain.ID, 1)

		block := chain.NewBlock(proof, prevHash)
		response := models.MineResponse{
			Message: "A new block was mined.",
			Index: block.Index,
			Transactions: block.Transactions,
			Proof: block.Proof,
			PreviousHash: block.PreviousHash,
		}

		if err := unmarshalToResponse(w, response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[MineHandler] - Failed to unmarshal GetMine response", err)
		}		
	}
}

func ConsensusHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ConsensusHandler] - Received request")

		replaced := chain.ResolveConflicts()
		response := models.ConsensusResponse{
			Message: "",
			Chain: chain.Chain,
		}

		if replaced {
			response.Message = fmt.Sprintf("Chain on node %s was replaced", chain.ID)
		} else {
			response.Message = fmt.Sprintf("Chain on node %s is authoritative", chain.ID)
		}

		if err := unmarshalToResponse(w, response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("[ConsensusHandler] - Failed to unmarshal ConsensusResponse", err)
		}
	}
}

func NewTransactionHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[NewTransactionHandler] - Received request")

		decoder := json.NewDecoder(r.Body)
		var transaction models.Transaction

		if err := decoder.Decode(&transaction); err == nil {
			index := chain.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)
			response := models.MessageResponse{ Message: fmt.Sprintf("Transaction added to block at index %d", index) }

			if err := unmarshalToResponse(w, response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[NewTransactionHandler] - Failed to unmarshal PostNewTransaction response", err)
			}
		} else {
			http.Error(w, "Request body missing transaction fields: expected a sender, recipient, and famount.", http.StatusBadRequest)
			log.Println("[NewTransactionHandler] - Did not receive needed fields in POST body, failed to decode body to Transaction", err)
		}
		defer r.Body.Close()
	}
}

func RegisterNodesHandler(chain *blockchain.Blockchain) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[RegisterNodeHandler] - Received request")
		
		decoder := json.NewDecoder(r.Body)
		var reqBody models.RegisterNodesRequest
		
		if err := decoder.Decode(&reqBody); err == nil {
			numNodesToAdd := len(reqBody.URLs)
			if numNodesToAdd == 0 {
				http.Error(w, "Request body had an empty list of nodes.", http.StatusBadRequest)
				return
			}

			for _, url := range reqBody.URLs {
				chain.RegisterNode(url)
			}

			log.Printf("[RegisterNodeHandler] - Registered %d new nodes %v\n", numNodesToAdd, reqBody.URLs)

			response := models.RegisterNodesResponse{ 
				Message: fmt.Sprintf("Registered new nodes to node %s", chain.ID), 
				TotalNodes: len(chain.Nodes),
				Nodes: chain.GetNodes(),
			}
			
			if err := unmarshalToResponse(w, response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[RegisterNodesnHandler] - Failed to unmarshal RegisterNodesResponse", err)
			}	
		} else {
			http.Error(w, "Request body missing url field.", http.StatusBadRequest)
			log.Println("[RegisterNodeHandler] - Did not receive needed fields in POST body, failed to decode body to RegisterNodeBody", err)
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
