package llmsupplier

import (
	"bytes"
	"context"
	_ "embed"
	"vibesiter/pkg/dto"

	"github.com/cockroachdb/errors"
)

func (r *Supplier) DesignSite(ctx context.Context, slug string, userPrompt string) (dto.SiteDesign, error) {
	var (
		output dto.SiteDesign
		buf    bytes.Buffer
	)

	err := r.systemPromptTemplateDesignSite.Execute(&buf, map[string]any{"Slug": slug})
	if err != nil {
		return dto.SiteDesign{}, errors.Wrap(err, "execute design site")
	}

	err = r.chat(ctx, buf.String(), userPrompt, &output)
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
	var (
		output    dto.Files
		bufUser   bytes.Buffer
		bufSystem bytes.Buffer
	)

	err := r.systemPromptTemplateGenerate.Execute(&bufSystem, map[string]any{"Slug": meta.Slug})
	if err != nil {
		return dto.Files{}, errors.Wrap(err, "execute design site")
	}

	err = r.userPromptTemplateGenerate.Execute(
		&bufUser,
		map[string]any{"Slug": meta.Slug, "UserPrompt": userPrompt, "Meta": meta, "Design": design},
	)
	if err != nil {
		return dto.Files{}, errors.Wrap(err, "execute design site")
	}

	err = r.chat(ctx, bufSystem.String(), bufUser.String(), &output)
	if err != nil {
		return dto.Files{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
