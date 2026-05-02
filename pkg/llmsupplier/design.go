package llmsupplier

import (
	"context"
	"vibesiter/pkg/dto"

	"github.com/cockroachdb/errors"
)

func (r *Supplier) DesignSite(ctx context.Context, slug string, userPrompt string) (dto.SiteDesign, error) {
	systemPrompt, err := r.prompts.Render("design", map[string]any{"Slug": slug})
	if err != nil {
		return dto.SiteDesign{}, errors.Wrap(err, "execute design site")
	}

	var output dto.SiteDesign

	_, err = r.chat(ctx, systemPrompt, userPrompt, &output)
	if err != nil {
		return dto.SiteDesign{}, errors.Wrap(err, "chat")
	}

	return output, nil
}

func (r *Supplier) GenerateSite(
	ctx context.Context,
	userPrompt string,
	meta SiteMeta,
	design dto.SiteDesign,
) (dto.Files, error) {
	var output dto.Files

	systemPrompt, err := r.prompts.Render("generate_system", map[string]any{"Slug": meta.Slug})
	if err != nil {
		return dto.Files{}, errors.Wrap(err, "execute design site")
	}

	userPrompt, err = r.prompts.Render("generate_user",
		map[string]any{"Slug": meta.Slug, "UserPrompt": userPrompt, "Meta": meta, "Design": design},
	)
	if err != nil {
		return dto.Files{}, errors.Wrap(err, "execute design site")
	}

	_, err = r.chat(ctx, systemPrompt, userPrompt, &output)
	if err != nil {
		return dto.Files{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
