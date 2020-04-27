package jsonld

// ErrorMessage describes error message interface
type ErrorMessage interface {
	GetCode() uint32
	GetMessage() string
}

// Error is the error holder
type Error struct {
	StatusCode  int    `json:"statusCode,omitempty"`
	Code        string `json:"code,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewError returns a json-ld error holder
func NewError(context, id, code, title string) *Resource {
	return &Resource{
		Context: context,
		Type:    "Error",
		NodeID:  id,
		Body: &Error{
			Code:       code,
			Title:      title,
			StatusCode: 400,
		},
	}
}

// WrapError returns a json-ld error holder
func WrapError(context string, id string, code string, err ErrorMessage) *Resource {
	return &Resource{
		Context: context,
		Type:    "Error",
		NodeID:  id,
		Body: &Error{
			Code:       code,
			Title:      err.GetMessage(),
			StatusCode: int(err.GetCode()),
		},
	}
}
