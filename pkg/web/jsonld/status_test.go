package jsonld_test

import (
	"encoding/json"
	"testing"

	"github.com/scraly/go.common/pkg/web/jsonld"
	"github.com/stretchr/testify/require"
)

func TestStatusJSONMashaller(t *testing.T) {
	res := jsonld.NewStatus("http://schema.org", "http://entry.cloud/api/v1/directories", "API/ERR/0001", 400)

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://entry.cloud/api/v1/directories\",\"@type\":\"Status\",\"statusCode\":400,\"title\":\"API/ERR/0001\"}", string(out), "Encoded json should be as expected")
}
