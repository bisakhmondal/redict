package persistence

import (
	"bytes"
	"encoding/json"
	"io"
)


// Marshal is a function that marshals the object into an io.Reader.
var Marshal = func(v interface{})(io.Reader, error){
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
