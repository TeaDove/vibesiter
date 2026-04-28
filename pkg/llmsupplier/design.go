package llmsupplier

import (
	"bytes"
	"context"
	_ "embed"
	"vibesiter/pkg/dto"

	"github.com/cockroachdb/errors"
)

//go:embed design.gohtml
var systemPromptDesignSite string

//go:embed generate_system.gohtml
var systemPromptGenerate string

//go:embed generate_user.gohtml
var userPromptGenerate string

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

type File struct {
	Path    string `json:"path"    validate:"required,max=128"`
	Content string `json:"content" validate:"required,max=1000000"`
}

type Files struct {
	Files []File `json:"files" validate:"required,min=1,dive"`
}

func (r *Supplier) GenerateSite(
	ctx context.Context,
	userPrompt string,
	meta SiteMeta,
	design dto.SiteDesign,
) (Files, error) {
	var (
		output    Files
		bufUser   bytes.Buffer
		bufSystem bytes.Buffer
	)

	err := r.systemPromptTemplateGenerate.Execute(&bufSystem, map[string]any{"Slug": meta.Slug})
	if err != nil {
		return Files{}, errors.Wrap(err, "execute design site")
	}

	err = r.userPromptTemplateGenerate.Execute(
		&bufUser,
		map[string]any{"Slug": meta.Slug, "UserPrompt": userPrompt, "Meta": meta, "Design": design},
	)
	if err != nil {
		return Files{}, errors.Wrap(err, "execute design site")
	}

	err = r.chat(ctx, bufSystem.String(), bufUser.String(), &output)
	if err != nil {
		return Files{}, errors.Wrap(err, "chat")
	}

	return output, nil
}
