package kvrepo

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Set(ctx context.Context, appID uuid.UUID, key string, value any) error {
	kv := ApplicationKV{ApplicationID: appID, Key: key, Value: value}

	err := r.db.WithContext(ctx).Save(&kv).Error
	if err != nil {
		return errors.Wrap(err, "save kv")
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, appID uuid.UUID, key string) (any, error) {
	kv, err := gorm.G[ApplicationKV](r.db).
		Where("application_id = ?", appID).
		Where("key = ?", key).Take(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get kv")
	}

	return kv.Value, nil
}
