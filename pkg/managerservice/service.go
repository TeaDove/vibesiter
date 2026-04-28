package managerservice

import (
	"context"
	"time"
	"vibesiter/pkg/llmsupplier"
	"vibesiter/pkg/managerrepo"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type Service struct {
	managerRepo *managerrepo.Repo
	llmSupplier *llmsupplier.Supplier
}

func NewService(managerRepo *managerrepo.Repo, llmSupplier *llmsupplier.Supplier) *Service {
	return &Service{managerRepo, llmSupplier}
}

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

	err = r.managerRepo.SaveApplication(ctx, &application)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "save application")
	}

	files, err := r.llmSupplier.GenerateSite(ctx, userPrompt, meta, design)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "generate site")
	}

	err = r.managerRepo.SaveFiles(ctx, application.ID, files.Files)
	if err != nil {
		return managerrepo.Application{}, errors.Wrap(err, "save files")
	}

	return application, nil
}
