package llmsupplier

import (
	"context"
	"vibesiter/pkg/dto"
	"vibesiter/pkg/managerrepo"

	"github.com/cockroachdb/errors"
)

func (r *Supplier) HTTP(
	ctx context.Context,
	app *managerrepo.Application,
	req *dto.HTTPRequest,
) (dto.LLMHTTPResponse, error) {
	systemPrompt, err := r.prompts.Render("http_system", map[string]any{"Slug": app.Slug})
	if err != nil {
		return dto.LLMHTTPResponse{}, errors.Wrap(err, "execute design site")
	}

	userPrompt, err := r.prompts.Render("http_user",
		map[string]any{
			"Slug":       app.Slug,
			"UserPrompt": app.UserPrompt,
			"MetaJson":   map[string]any{"slug": app.Slug, "title": app.Title, "description": app.Description},
			"Design":     app.Design,
			"Request":    req,
		},
	)
	if err != nil {
		return dto.LLMHTTPResponse{}, errors.Wrap(err, "execute design site")
	}

	var output dto.LLMHTTPResponse

	err = r.chat(ctx, systemPrompt, userPrompt, &output)
	if err != nil {
		return dto.LLMHTTPResponse{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
