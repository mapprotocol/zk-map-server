package entity

// ============================== request ==============================

type GetProofRequest struct {
	ChainID uint16 `form:"chain_id"`
	Height  string `form:"height"`
}

// ============================== response ==============================

type GetProofResponse struct {
	Height   string `json:"height"`
	Status   uint8  `json:"status"`
	Result   Result `json:"result,omitempty"`
	ErrorMsg string `json:"error_msg"`
}

type Result struct {
	Proof       *Proof   `json:"proof,omitempty"`
	PublicInput []string `json:"public_input,omitempty"`
}

type Proof struct {
	PiA      []string   `json:"pi_a"`
	PiB      [][]string `json:"pi_b"`
	PiC      []string   `json:"pi_c"`
	Protocol string     `json:"protocol"`
}
