package jsonld

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var (
	jsonpbMarshaler = &jsonpb.Marshaler{OrigName: true}
)

// Resource represents the JSON-LD header
type Resource struct {
	Context string      `json:"@context,omitempty"`
	NodeID  string      `json:"@id,omitempty"`
	Type    string      `json:"@type,omitempty"`
	Body    interface{} `json:",omitempty"`
}

// NewResource returns a JSONLD resource
func NewResource(context, id, _type string, body interface{}) *Resource {
	return &Resource{
		Context: context,
		NodeID:  id,
		Type:    _type,
		Body:    body,
	}
}

// -----------------------------------------------------------------------------

// MarshalJSON is used to export resource as a JSON encoded payload
func (r Resource) MarshalJSON() ([]byte, error) {
	out := fmt.Sprintf(`{"@context":"%s","@id":"%s","@type":"%s",`, r.Context, r.NodeID, r.Type)

	var body string
	var err error

	// Process body
	switch r.Body.(type) {
	case proto.Message:
		msg := r.Body.(proto.Message)
		body, err = jsonpbMarshaler.MarshalToString(msg)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to marshal protobuf message")
		}

	default:
		var out []byte
		out, err = json.Marshal(r.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to marshal generic payload")
		}
		body = string(out)
	}

	// Excludes '{' and '}'
	out += body[1 : len(body)-1]
	out += "}"
	return []byte(out), nil
}
