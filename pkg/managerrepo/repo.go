package managerrepo

import (
	"context"

	"github.com/cockroachdb/errors"
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
		return errors.Wrap(err, "save message")
	}

	zerolog.Ctx(ctx).Info().
		Object("msg", v).
		Msg("saved")

	return nil
}
