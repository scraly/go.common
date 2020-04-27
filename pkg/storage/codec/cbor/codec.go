/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package cbor

import (
	"bytes"

	api "github.com/scraly/go.common/pkg/storage/codec"

	"github.com/ugorji/go/codec"
)

type cborCodec struct {
	mh codec.Handle
}

func (j cborCodec) Marshal(v interface{}) ([]byte, error) {
	// Pepare encoder
	w := new(bytes.Buffer)
	err := codec.NewEncoder(w, j.mh).Encode(v)
	// Return no-error
	return w.Bytes(), err
}

func (j cborCodec) Unmarshal(d []byte, v interface{}) error {
	// Pepare decoder
	rdr := bytes.NewReader(d)

	// Return no-error
	return codec.NewDecoder(rdr, j.mh).Decode(v)
}

func (j cborCodec) String() string {
	return "cbor"
}

// NewCodec returns a JSON codec
func NewCodec() api.Codec {
	return cborCodec{
		mh: &codec.CborHandle{},
	}
}
