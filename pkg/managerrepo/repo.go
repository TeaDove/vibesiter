package managerrepo

import (
	"context"
	"vibesiter/pkg/llmsupplier"

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

func (r *Repo) SaveApplication(ctx context.Context, v *Application) error {
	err := r.db.WithContext(ctx).Save(v).Error
	if err != nil {
		return errors.Wrap(err, "save app")
	}

	zerolog.Ctx(ctx).Info().
		Object("msg", v).
		Msg("saved")

	return nil
}

func (r *Repo) SaveFiles(ctx context.Context, appId uuid.UUID, files []llmsupplier.File) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, file := range files {
			err := tx.WithContext(ctx).Save(&ApplicationFiles{
				ApplicationID: appId,
				Path:          file.Path,
				Content:       file.Content,
			}).Error
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
