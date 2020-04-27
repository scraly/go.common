/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package msgpack

import (
	"bytes"

	api "github.com/scraly/go.common/pkg/storage/codec"

	"github.com/ugorji/go/codec"
)

type msgpackCodec struct {
	mh codec.Handle
}

func (j msgpackCodec) Marshal(v interface{}) ([]byte, error) {
	// Pepare encoder
	w := new(bytes.Buffer)
	err := codec.NewEncoder(w, j.mh).Encode(v)
	// Return no-error
	return w.Bytes(), err
}

func (j msgpackCodec) Unmarshal(d []byte, v interface{}) error {
	// Pepare decoder
	rdr := bytes.NewReader(d)

	// Return no-error
	return codec.NewDecoder(rdr, j.mh).Decode(v)
}

func (j msgpackCodec) String() string {
	return "msgpack"
}

// NewCodec returns a JSON codec
func NewCodec() api.Codec {
	return msgpackCodec{
		mh: &codec.MsgpackHandle{RawToString: true, WriteExt: true},
	}
}
