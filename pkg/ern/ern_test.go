package ern

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseERN(t *testing.T) {
	cases := []struct {
		input string
		ern   ERN
		err   error
	}{
		{
			input: "invalid",
			err:   errors.New(invalidPrefix),
		},
		{
			input: "ern:nope",
			err:   errors.New(invalidSections),
		},
		{
			input: "ern:12:user:123456789",
			ern: ERN{
				Tenant:   "12",
				Type:     "user",
				Resource: "123456789",
			},
		}, {
			input: "ern::group:entry/administrators",
			ern: ERN{
				Tenant:   "",
				Type:     "group",
				Resource: "entry/administrators",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			spec, err := Parse(tc.input)
			if tc.ern != spec {
				t.Errorf("Expected %q to parse as %v, but got %v", tc.input, tc.ern, spec)
			}
			if err == nil && tc.err != nil {
				t.Errorf("Expected err to be %v, but got nil", tc.err)
			} else if err != nil && tc.err == nil {
				t.Errorf("Expected err to be nil, but got %v", err)
			} else if err != nil && tc.err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected err to be %v, but got %v", tc.err, err)
			}
		})
	}
}

func TestERNString(t *testing.T) {
	res := ERN{
		Tenant:   "12",
		Type:     "user",
		Resource: "123456789",
	}
	require.Equal(t, res.String(), "ern:12:user:123456789")
}
