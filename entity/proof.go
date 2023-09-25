package entity

// ============================== request ==============================

type GetProofRequest struct {
	Height string `form:"height"`
}

// ============================== response ==============================

type GetProofResponse struct {
	Id       string `json:"id"`
	Status   uint8  `json:"status"`
	Result   Result `json:"result"`
	ErrorMsg string `json:"error_msg"`
}

type Result struct {
	Proof struct {
		PiA      []string   `json:"pi_a"`
		PiB      [][]string `json:"pi_b"`
		PiC      []string   `json:"pi_c"`
		Protocol string     `json:"protocol"`
	} `json:"proof"`
	PublicInput []string `json:"public_input"`
}
