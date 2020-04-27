/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package json

import (
	"bytes"

	api "github.com/scraly/go.common/pkg/storage/codec"

	"github.com/ugorji/go/codec"
)

type jsonCodec struct {
	mh codec.Handle
}

func (j jsonCodec) Marshal(v interface{}) ([]byte, error) {
	// Pepare encoder
	w := new(bytes.Buffer)
	err := codec.NewEncoder(w, j.mh).Encode(v)

	// Return
	return w.Bytes(), err
}

func (j jsonCodec) Unmarshal(d []byte, v interface{}) error {
	// Pepare decoder
	rdr := bytes.NewReader(d)
	return codec.NewDecoder(rdr, j.mh).Decode(v)
}

func (j jsonCodec) String() string {
	return "json"
}

// NewCodec returns a JSON codec
func NewCodec() api.Codec {
	return jsonCodec{
		mh: &codec.JsonHandle{},
	}
}
