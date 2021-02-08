package model

type Heos struct {
	Command string `json:"command"`
	Result  string `json:"result"`
	Message string `json:"message"`
}
