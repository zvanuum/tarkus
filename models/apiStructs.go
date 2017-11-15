package models

type MessageResponse struct {
	Message string `json:"message"`
}

type ChainResponse struct {
	Chain []Block `json:"chain"`
	Length int `json:"length"`
}

type MineResponse struct {
	Message string `json:"message"`
	Index int64 `json:"index"`
	Transactions []Transaction `json:"transactions"`
	Proof int64 `json:"proof"`
	PreviousHash string `json:"previousHash"`
}

type RegisterNodesRequest struct {
	URLs []string `json:"urls"`
}

type RegisterNodesResponse struct {
	Message string `json:"message"`
	TotalNodes int `json:"totalNodes"`
	Nodes []string `json:"nodes"`
}

type ConsensusResponse struct {
	Message string `json:"message"`
	Chain []Block `json:"chain"`
}
