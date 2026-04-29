package managerrepo

import (
	"context"
	"vibesiter/pkg/dto"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) InsertApplication(ctx context.Context, v *Application) error {
	err := r.db.WithContext(ctx).Create(v).Error
	if err != nil {
		return errors.Wrap(err, "save app")
	}

	zerolog.Ctx(ctx).Info().
		Object("msg", v).
		Msg("saved")

	return nil
}

func (r *Repo) InsertFiles(ctx context.Context, appId uuid.UUID, files []dto.File) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, file := range files {
			err := tx.WithContext(ctx).Save(
				&ApplicationFile{
					ApplicationID: appId,
					Path:          file.Path,
					Content:       []byte(file.Content),
				},
			).Error
			if err != nil {
				return errors.Wrap(err, "save files")
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "save files")
	}

	zerolog.Ctx(ctx).Info().
		Int("len", len(files)).
		Msg("files.saved")

	return nil
}

func (r *Repo) SelectApplicationBySlug(ctx context.Context, slug string) (Application, error) {
	v, err := gorm.G[Application](r.db).Where("slug = ?", slug).Take(ctx)
	if err != nil {
		return Application{}, errors.Wrap(err, "select by slug")
	}

	return v, nil
}

func (r *Repo) SelectFile(ctx context.Context, appID uuid.UUID, path string) (ApplicationFile, error) {
	v, err := gorm.G[ApplicationFile](r.db).
		Where("application_id = ?", appID.String()).
		Where("path = ?", path).
		Take(ctx)
	if err != nil {
		return ApplicationFile{}, errors.Wrap(err, "select by app and path")
	}

	return v, nil
}
