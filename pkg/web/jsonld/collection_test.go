package jsonld_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/require"

	"github.com/scraly/go.common/pkg/web/jsonld"
	"github.com/scraly/go.common/pkg/web/paginator"
)

type jsonLdTest struct {
	Ern                string               `protobuf:"bytes,4,opt,name=ern" json:"ern,omitempty"`
	PasswordModifiedAt *timestamp.Timestamp `protobuf:"bytes,8,opt,name=password_modified_at,json=passwordModifiedAt" json:"password_modified_at,omitempty"`
}

func TestCollectionJSONMashallerWhenNotASlice(t *testing.T) {
	res := jsonld.NewCollection("http://schema.org", "http://toto.continental.cloud/api/v1/directories", 8)

	_, err := json.Marshal(res)
	require.Error(t, err, "Error should be raised on json marshalling")
}

func TestCollectionJSONMashaller(t *testing.T) {
	res := jsonld.NewCollection("http://schema.org", "http://toto.continental.cloud/api/v1/directories", []interface{}{})

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://toto.continental.cloud/api/v1/directories\",\"@type\":\"Collection\",\"total\":\"0\",\"members\":[]}", string(out), "Encoded json should be as expected")
}

func TestCollectionJSONMashallerWithProto(t *testing.T) {
	modified, err := ptypes.TimestampProto(time.Date(2018, time.March, 13, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err, "Error should be nil")

	collection := []*jsonLdTest{
		{
			Ern:                "123456",
			PasswordModifiedAt: modified,
		},
	}

	res := jsonld.NewCollection("http://schema.org", "http://toto.continental.cloud/api/v1/directories", collection)

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://toto.continental.cloud/api/v1/directories\",\"@type\":\"Collection\",\"total\":\"0\",\"members\":[{\"ern\":\"123456\",\"password_modified_at\":{\"seconds\":1520899200}}]}", string(out), "Encoded json should be as expected")
}

func TestCollectionJSONMashallerWithProtoAndPaginator(t *testing.T) {
	modified, err := ptypes.TimestampProto(time.Date(2018, time.March, 13, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err, "Error should be nil")

	collection := []*jsonLdTest{
		{
			Ern:                "123456",
			PasswordModifiedAt: modified,
		},
	}

	res := jsonld.NewCollection("http://schema.org", "http://toto.continental.cloud/api/v1/directories", collection)
	paginator := paginator.NewPaginator(3, 20)
	paginator.SetTotal(80)
	req, err := http.NewRequest("GET", "http://toto.continental.cloud/api/v1/directories", nil)
	require.NoError(t, err, "Error should not be raised on get directories")
	res.SetPaginator(req, paginator)

	out, err := json.Marshal(res)
	require.NoError(t, err, "Error should not be raised on json marshalling")
	require.Equal(t, "{\"@context\":\"http://schema.org\",\"@id\":\"http://toto.continental.cloud/api/v1/directories\",\"@type\":\"PagedCollection\",\"total\":\"80\",\"current_page\":\"3\",\"item_per_page\":\"20\",\"first_page\":\"http://toto.continental.cloud/api/v1/directories?page=1\",\"next_page\":\"http://toto.continental.cloud/api/v1/directories?page=4\",\"previous_page\":\"http://toto.continental.cloud/api/v1/directories?page=2\",\"last_page\":\"http://toto.continental.cloud/api/v1/directories?page=4\",\"members\":[{\"ern\":\"123456\",\"password_modified_at\":{\"seconds\":1520899200}}]}", string(out), "Encoded json should be as expected")
}
