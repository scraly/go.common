/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package json_test

import (
	"testing"

	"github.com/scraly/go.common/pkg/storage/codec/json"
	"github.com/stretchr/testify/require"
)

type Serializable struct {
	Property string `codec:"property"`
	Bool     bool   `codec:"bool"`
}

var (
	s = []Serializable{
		{
			Property: "toto",
			Bool:     true,
		},
		{
			Property: "tutu",
			Bool:     false,
		},
	}
)

func TestEncodeDecode(t *testing.T) {
	c := json.NewCodec()
	require.Equal(t, "json", c.String(), "Codec name should be 'json'")

	payload, err := c.Marshal(s)
	require.NoError(t, err, "Encoding should not raise error")
	require.NotNil(t, payload, "Payload should nt be nil")

	var result []Serializable
	err = c.Unmarshal(payload, &result)
	require.NoError(t, err, "Payload decoding shouyld not raise error")
}

func BenchmarkMarshal(b *testing.B) {
	c := json.NewCodec()
	for n := 0; n < b.N; n++ {
		_, err := c.Marshal(s)
		require.NoError(b, err, "Payload encoding should not raise error")
	}
}
