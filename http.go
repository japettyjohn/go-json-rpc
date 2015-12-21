package jsonrpc

import (
	"encoding/json"
	"net/http"
)

func MustJSONErrorHTTP(w http.ResponseWriter, errCode int, errMsg string, data interface{}) {
	if e := jsonEncodeHTTP(JSONMessage{Error: &JSONError{errCode, errMsg, data}}); e != nil {
		panic(e)
	}
}

func MustJSONResultHTTP(w http.ResponseWriter, d interface{}) {
	if d == nil {
		panic("MustJsonResult does not accept nil data.")
	}

	if e := jsonEncodeHTTP(JSONMessage{Result: d}); e != nil {
		panic(e)
	}
}

func jsonEncodeHTTP(w http.ResponseWriter, d interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(d)
}
