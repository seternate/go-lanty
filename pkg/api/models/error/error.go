package errormodel

type Error struct {
	Error ErrorDescriptor `json:"error"`
}

type ErrorDescriptor struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
