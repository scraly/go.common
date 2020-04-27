package jsonld

// Status is the response holder
type Status struct {
	StatusCode  int    `json:"statusCode,omitempty"`
	Code        string `json:"code,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewStatus returns a status
func NewStatus(context, id, title string, status int) *Resource {
	return &Resource{
		Context: context,
		Type:    "Status",
		NodeID:  id,
		Body: &Status{
			Title:      title,
			StatusCode: status,
		},
	}
}
