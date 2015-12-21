package jsonrpc

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (j JSONError) Error() string {
	return j.Message
}

type JSONMessage struct {
	Error  *JSONError  `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

type JSONMessageSimpleError struct {
	Error  *string     `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}
