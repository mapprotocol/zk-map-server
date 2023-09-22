package entity

// ============================== request ==============================

type ExampleRequest struct {
	Msg string `form:"msg"`
}

// ============================== response ==============================

type ExampleResponse struct {
	Msg string `json:"msg"`
}
