package model

type ErrorResponse struct {
	Code        string `json:"code"`
	HTTPCode    int    `json:"http_code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Internal    string `json:"internal"`
}
