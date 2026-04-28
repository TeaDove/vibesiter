package llmsupplier

import (
	"context"

	"github.com/cockroachdb/errors"
)

type HTTPRequest struct {
	Headers map[string]string `json:"headers"`
	Body    any               `json:"body"`
}

type HTTPResponse struct {
	StatusCode  int    `json:"statusCode"  validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
	Body        any    `json:"body"`
}

func (r *Supplier) HTTP(ctx context.Context, backendPrompt string, req *HTTPRequest) (HTTPResponse, error) {
	var output HTTPResponse

	err := r.chat(ctx, "", "", output)
	if err != nil {
		return HTTPResponse{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
