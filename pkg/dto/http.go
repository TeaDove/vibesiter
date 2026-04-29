package dto

type HTTPRequest struct {
	Method string `json:"method"         validate:"required"`
	Path   string `json:"path"           validate:"required"`
	Body   any    `json:"body,omitempty"`
}

type HTTPResponse struct {
	Status      int    `json:"status"                validate:"required"`
	ContentType string `json:"contentType,omitempty"`
	Body        any    `json:"body,omitempty"`
}
