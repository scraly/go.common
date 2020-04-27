package jsonld_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"

	"github.com/scraly/go.common/pkg/web/jsonld"
)

func TestResourceJSONMashaller(t *testing.T) {
	res := jsonld.NewResource("http://schema.org", "http://toto.continental.cloud/api/v1/directories", "Collection", map[string]interface{}{
		"perPage": 20,
	})

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://toto.continental.cloud/api/v1/directories\",\"@type\":\"Collection\",\"perPage\":20}", string(out), "Encoded json should be as expected")
}

func TestResourceJSONMashallerWithProto(t *testing.T) {
	modified, err := ptypes.TimestampProto(time.Date(2018, time.March, 13, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err, "Error should be nil")

	proto := &jsonLdTest{
		Ern:                "123456",
		PasswordModifiedAt: modified,
	}
	res := jsonld.NewResource("http://schema.org", "http://toto.continental.cloud/api/v1/directories", "Account", proto)

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://toto.continental.cloud/api/v1/directories\",\"@type\":\"Account\",\"ern\":\"123456\",\"password_modified_at\":{\"seconds\":1520899200}}", string(out), "Encoded json should be as expected")
}
