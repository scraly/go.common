package jsonld

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/scraly/go.common/pkg/web/paginator"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// CollectionResource represents the JSON-LD header
type CollectionResource struct {
	Resource

	Total        uint        `json:"total"`
	ItemPerPage  uint        `json:"itemPerPage,omitempty"`
	CurrentPage  uint        `json:"currentPage,omitempty"`
	FirstPage    string      `json:"firstPage,omitempty"`
	NextPage     string      `json:"nextPage,omitempty"`
	PreviousPage string      `json:"previousPage,omitempty"`
	LastPage     string      `json:"lastPage,omitempty"`
	Members      interface{} `json:"members"`
}

// NewCollection returns a JSONLD Collection resource
func NewCollection(context, id string, members interface{}) *CollectionResource {
	return &CollectionResource{
		Resource: Resource{
			Context: context,
			NodeID:  id,
			Type:    "Collection",
		},
		Members: members,
	}
}

// SetPaginator defines values of the JSONLD Collection according to pagination request.
func (j *CollectionResource) SetPaginator(r *http.Request, paginator *paginator.Pagination) {

	// Check pagination usage
	if paginator == nil {
		return
	}

	j.Total = paginator.Total()

	if paginator.HasOtherPages() {
		j.Type = "PagedCollection"
		j.ItemPerPage = paginator.PerPage
		j.CurrentPage = paginator.Page
	}

	q := r.URL.Query()

	if paginator.HasOtherPages() {
		q.Set("page", fmt.Sprintf("%d", 1))
		r.URL.RawQuery = q.Encode()
		j.FirstPage = r.URL.String()
	}
	if paginator.HasPrev() {
		q.Set("page", fmt.Sprintf("%d", paginator.Page-1))
		r.URL.RawQuery = q.Encode()
		j.PreviousPage = r.URL.String()
	}
	if paginator.HasNext() {
		q.Set("page", fmt.Sprintf("%d", paginator.Page+1))
		r.URL.RawQuery = q.Encode()
		j.NextPage = r.URL.String()
	}
	if paginator.HasOtherPages() {
		q.Set("page", fmt.Sprintf("%d", paginator.NumPages()))
		r.URL.RawQuery = q.Encode()
		j.LastPage = r.URL.String()
	}

}

// -----------------------------------------------------------------------------

var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

// MarshalJSON is used to export resource as a JSON encoded payload
func (j CollectionResource) MarshalJSON() ([]byte, error) {
	// Preconditions
	v := reflect.ValueOf(j.Members)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Members must be a slice")
	}

	out := fmt.Sprintf(`{"@context":"%s","@id":"%s","@type":"%s", "total":"%d"`, j.Context, j.NodeID, j.Type, j.Total)

	var elements []string
	if j.CurrentPage > 0 {
		elements = append(elements, fmt.Sprintf(`"current_page":"%d"`, j.CurrentPage))
	}
	if j.ItemPerPage > 0 {
		elements = append(elements, fmt.Sprintf(`"item_per_page":"%d"`, j.ItemPerPage))
	}
	if j.FirstPage != "" {
		elements = append(elements, fmt.Sprintf(`"first_page":"%s"`, j.FirstPage))
	}
	if j.NextPage != "" {
		elements = append(elements, fmt.Sprintf(`"next_page":"%s"`, j.NextPage))
	}
	if j.PreviousPage != "" {
		elements = append(elements, fmt.Sprintf(`"previous_page":"%s"`, j.PreviousPage))
	}
	if j.LastPage != "" {
		elements = append(elements, fmt.Sprintf(`"last_page":"%s"`, j.LastPage))
	}

	if len(elements) > 0 {
		out = fmt.Sprintf(`%s, %s`, out, strings.Join(elements, ","))
	}

	var members string

	collectionType := reflect.ValueOf(j.Members)
	result := reflect.New(reflect.TypeOf(j.Members).Elem())

	if result.Elem().Type().AssignableTo(protoMessageType) {
		var messages []string
		for i := 0; i < collectionType.Len(); i++ {
			memberBody, err := jsonpbMarshaler.MarshalToString(collectionType.Index(i).Interface().(proto.Message))
			if err != nil {
				return nil, errors.Wrap(err, "Unable to marshal protobuf collection")
			}
			messages = append(messages, memberBody)
		}

		members = fmt.Sprintf(`[%s]`, strings.Join(messages, ","))
	} else {
		body, err := json.Marshal(j.Members)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to marshal generic collection")
		}
		members = string(body)
	}

	out = fmt.Sprintf(`%s, "members":%s}`, out, members)

	return []byte(out), nil
}
