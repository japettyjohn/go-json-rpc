package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	// Protocol version that must be specified on every RPC message
	ProtocolVersion = "2.0"
	// ParseErrorCode for invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ParseErrorCode = -32700
	// InvalidRequestCode for JSON is not a valid Request object.
	InvalidRequestCode = -32600
	// MethodNotFoundCode when method does not exist / is not available.
	MethodNotFoundCode = -32601
	// InvalidParamsCode used with invalid method parameter(s).
	InvalidParamsCode = -32602
	// InternalErrorCode to show internal JSON-RPC error.
	InternalErrorCode = -32603

	// The error codes from and including -32768 to -32000 are reserved for pre-defined errors
	// Server error = -32000 to -32099 	 	Reserved for implementation-defined server-errors.
)

// Error implements error on responses, see http://www.jsonrpc.org/specification#error_object
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (j Error) Error() string {
	return j.Message
}

// Resposne is what is sent in reply to any non-notification requests, see http://www.jsonrpc.org/specification#response_object
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Error   *Error      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	ID      interface{} `json:"id"`
}

func (r Response) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

// Request is the form all incoming requests are parsed into, see http://www.jsonrpc.org/specification#request_object
type Request struct {
	JSONRPC  string      `json:"jsonrpc"`
	Method   string      `json:"method"`
	Params   interface{} `json:"params,omitempty"`
	ID       interface{} `json:"id,omitempty"`
	Response Response    `json:"-"`
}

// DecodeRequest parses out a single request, params are parsed into a generic type
func DecodeRequest(d *json.Decoder) *Request {
	return DecodeRequestTypedParams(d, nil)
}

// DecodeRequestTypedParams parses out a single request, params are parsed into the provided pointer
func DecodeRequestTypedParams(d *json.Decoder, params interface{}) *Request {
	r := &Request{Response: Response{JSONRPC: ProtocolVersion}, Params: params}
	if e := d.Decode(r); e != nil {
		r.Response.Error = &Error{Code: ParseErrorCode, Message: e.Error()}
		return r
	}

	r.Response.ID = r.ID
	if r.Method == "" {
		r.Response.Error = &Error{Code: InvalidRequestCode, Message: "Invalid method: ''."}
	} else if r.JSONRPC != ProtocolVersion {
		r.Response.Error = &Error{Code: InvalidRequestCode, Message: fmt.Sprintf("Invalid jsonrpc: '%s'.", r.JSONRPC)}
	}
	return r
}
