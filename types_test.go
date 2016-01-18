package jsonrpc

import (
	"bytes"
	"encoding/json"
	"testing"
)

type paramsType struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p paramsType) String() string {
	return p.FirstName + " " + p.LastName
}

func TestDecoding(t *testing.T) {
	call := bytes.NewReader([]byte(`{"jsonrpc":"2.0","method":"test", "params": {"firstName":"Joe","lastName":"D."}, "id":1}`))
	req := DecodeRequest(json.NewDecoder(call))
	if req.Response.Error != nil {
		t.Fatal(req.Response.Error)
	}

	//	 Check ID
	if id, ok := req.Response.ID.(float64); !ok {
		t.Fatalf("Did not get expected id type of float64")
	} else if id != float64(1) {
		t.Fatalf("Did not get expected id of 1.")
	}

	// Check protocol
	if req.Response.JSONRPC != ProtocolVersion {
		t.Fatalf("Unexpected protocol '%s'", req.Response.JSONRPC)
	}

	// Check method
	if req.Method != "test" {
		t.Fatalf("Unexpected method '%s'", req.Method)
	}

	// Provide the params
	call.Seek(0, 0)
	req = DecodeRequestTypedParams(json.NewDecoder(call), &paramsType{})

	if p, ok := req.Params.(*paramsType); !ok {
		t.Fatalf("Unexpected paramsType %T", req.Params)
	} else if n, expected := p.String(), "Joe D."; n != expected {
		t.Fatalf("Expected name '%s' but got '%s'.", expected, n)
	}

}
