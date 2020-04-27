package jsonld_test

import (
	"encoding/json"
	"testing"

	"github.com/scraly/go.common/pkg/web/jsonld"
	"github.com/stretchr/testify/require"
)

func TestErrorJSONMashaller(t *testing.T) {
	res := jsonld.NewError("http://schema.org", "http://entry.cloud/api/v1/directories", "API/ERR/0001", "Unable to process your request")

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://entry.cloud/api/v1/directories\",\"@type\":\"Error\",\"statusCode\":400,\"code\":\"API/ERR/0001\",\"title\":\"Unable to process your request\"}", string(out), "Encoded json should be as expected")
}
