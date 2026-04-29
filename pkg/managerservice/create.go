package managerservice

import (
	"context"
	"fmt"
	"time"
	"vibesiter/pkg/managerrepo"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/teadove/teasutils/utils/random_utils"
	"gorm.io/gorm"
)

func (r *Service) GenerateApp(ctx context.Context, userPrompt string) (managerrepo.Application, error) {
	if userPrompt == "" {
		return managerrepo.Application{}, errors.New("userPrompt is required")
	}

	now := time.Now()
	application := managerrepo.Application{
		ID:             uuid.Must(uuid.NewV7()),
		CreatedAt:      now,
		UpdatedAt:      now,
		LastAccessedAt: now,
		Status:         managerrepo.ApplicationStatusCREATING,
		UserPrompt:     userPrompt,
	}

	validationResponse, err := r.llmSupplier.ValidateVibeSite(ctx, userPrompt)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "validate vibe site")
	}

	if validationResponse.Success <= 0.2 {
		return managerrepo.Application{}, errors.Newf("validation site error: %s", validationResponse.Description)
	}

	meta, err := r.llmSupplier.ExtractSiteMeta(ctx, userPrompt)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "extract site meta")
	}

	design, err := r.llmSupplier.DesignSite(ctx, meta.Slug, userPrompt)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "design site")
	}

	application.Title = meta.Title
	application.Description = meta.Description
	application.Slug = meta.Slug
	application.Design = design

	err = r.createApplication(ctx, &application)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "create application")
	}

	meta.Slug = application.Slug

	files, err := r.llmSupplier.GenerateSite(ctx, userPrompt, meta, design)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "generate site")
	}

	err = r.managerRepo.InsertFiles(ctx, application.ID, files.Files)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "save files")
	}

	return application, nil
}

func (r *Service) createApplication(ctx context.Context, application *managerrepo.Application) error {
	originalSlug := application.Slug

	const maxLen = 10

	text := random_utils.TextWithLen(maxLen)

	for i := range maxLen {
		err := r.managerRepo.InsertApplication(ctx, application)
		if err == nil {
			return nil
		}

		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.Wrap(err, "save application")
		}

		application.Slug = fmt.Sprintf("%s-%s", originalSlug, text[:i])
	}

	return errors.New("failed to create slug")
}

func (r *Service) ListApps(ctx context.Context) ([]managerrepo.Application, error) {
	return r.managerRepo.SelectApplications(ctx)
}
