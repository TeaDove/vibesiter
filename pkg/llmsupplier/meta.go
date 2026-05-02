package llmsupplier

import (
	"context"

	"github.com/cockroachdb/errors"
)

type ValidationResponse struct {
	Success     float32 `json:"success"               validate:"required"`
	Description string  `json:"description,omitempty"`
}

func (r *Supplier) ValidateVibeSite(ctx context.Context, userPrompt string) (ValidationResponse, error) {
	sysPrompt, err := r.prompts.Render("meta_validate", nil)
	if err != nil {
		return ValidationResponse{}, errors.Wrap(err, "render")
	}

	var output ValidationResponse

	_, err = r.chat(ctx, sysPrompt, userPrompt, &output)
	if err != nil {
		return ValidationResponse{}, errors.Wrap(err, "chat")
	}

	return output, nil
}

type SiteMeta struct {
	Title       string `json:"title"       validate:"required"`
	Description string `json:"description" validate:"required"`
	Slug        string `json:"slug"        validate:"required"`
}

func (r *Supplier) ExtractSiteMeta(ctx context.Context, userPrompt string) (SiteMeta, error) {
	var output SiteMeta

	sysPrompt, err := r.prompts.Render("meta", nil)
	if err != nil {
		return SiteMeta{}, errors.Wrap(err, "render")
	}

	_, err = r.chat(ctx, sysPrompt, userPrompt, &output)
	if err != nil {
		return SiteMeta{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
