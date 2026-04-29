package llmsupplier

import (
	"bytes"
	"context"
	"vibesiter/pkg/dto"
	"vibesiter/pkg/managerrepo"

	"github.com/cockroachdb/errors"
)

func (r *Supplier) HTTP(
	ctx context.Context,
	app *managerrepo.Application,
	req *dto.HTTPRequest,
) (dto.HTTPResponse, error) {
	var (
		output    dto.HTTPResponse
		bufUser   bytes.Buffer
		bufSystem bytes.Buffer
	)

	err := r.systemPromptTemplateHTTP.Execute(&bufSystem, map[string]any{"Slug": app.Slug})
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "execute design site")
	}

	err = r.userPromptTemplateHTTP.Execute(&bufUser,
		map[string]any{
			"Slug":       app.Slug,
			"UserPrompt": app.UserPrompt,
			"MetaJson":   map[string]any{"slug": app.Slug, "title": app.Title, "description": app.Description},
			"Design":     app.Design,
			"Request":    req,
		},
	)
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "execute design site")
	}

	err = r.chat(ctx, bufSystem.String(), bufUser.String(), &output)
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
