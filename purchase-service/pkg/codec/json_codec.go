// Package codec registers a JSON codec for gRPC so plain Go structs can be used
// instead of protobuf-generated types. This is used because we manage proto files
// in a separate repository and generate code via GitHub Actions (contract-first approach).
package codec

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

// JSONCodec is a gRPC codec that uses JSON encoding.
type JSONCodec struct{}

func (JSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) Name() string {
	return "proto"
}
