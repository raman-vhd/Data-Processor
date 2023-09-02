package model

type APIResponse struct {
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}
