package dtos

type Request struct {
	Question string `json:"question"`
}

type Response struct {
	Answer string `json:"answer"`
}
