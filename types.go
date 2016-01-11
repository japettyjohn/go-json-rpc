package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	ProtocolVersion = "2.0"
	// 	Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ParseErrorCode = -32700
	// The JSON sent is not a valid Request object.
	InvalidRequestCode = -32600
	// The method does not exist / is not available.
	MethodNotFoundCode = -32601
	// Invalid method parameter(s).
	InvalidParamsCode = -32602
	// Internal JSON-RPC error.
	InternalErrorCode = -32603

	// The error codes from and including -32768 to -32000 are reserved for pre-defined errors
	// Server error = -32000 to -32099 	 	Reserved for implementation-defined server-errors.
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (j Error) Error() string {
	return j.Message
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Error   *Error      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Id      interface{} `json:"id"`
}

func (r Response) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

type Request struct {
	JSONRPC  string      `json:"jsonrpc"`
	Method   string      `json:"method"`
	Params   interface{} `json:"params,omitempty"`
	Id       interface{} `json:"id,omitempty"`
	Response Response    `json:"-"`
}

func DecodeRequest(d *json.Decoder) *Request {
	r := &Request{Response: Response{JSONRPC: ProtocolVersion}}
	if e := d.Decode(r); e != nil {
		r.Response.Error = &Error{Code: ParseErrorCode, Message: e.Error()}
		return r
	}

	r.Response.Id = r.Id
	if r.Method == "" {
		r.Response.Error = &Error{Code: InvalidRequestCode, Message: "Invalid method: ''."}
	} else if r.JSONRPC != ProtocolVersion {
		r.Response.Error = &Error{Code: InvalidRequestCode, Message: fmt.Sprintf("Invalid jsonrpc: '%s'.", r.JSONRPC)}
	}
	return r
}
