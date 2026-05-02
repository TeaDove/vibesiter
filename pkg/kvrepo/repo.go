package kvrepo

import (
	"context"
	"encoding/json"

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
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "marshal value")
	}

	kv := ApplicationKV{ApplicationID: appID, Key: key, Value: string(valueBytes)}

	err = r.db.WithContext(ctx).Save(&kv).Error
	if err != nil {
		return errors.Wrap(err, "save kv")
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, appID uuid.UUID, key string) (ApplicationKV, error) {
	kv, err := gorm.G[ApplicationKV](r.db).
		Where("application_id = ?", appID).
		Where("key = ?", key).Take(ctx)
	if err != nil {
		return ApplicationKV{}, errors.Wrap(err, "get kv")
	}

	return kv, nil
}

func (r *Repo) List(ctx context.Context, appID uuid.UUID, keyGlob string) ([]ApplicationKV, error) {
	kvs, err := gorm.G[ApplicationKV](r.db).
		Where("application_id = ?", appID).
		Where("key like ?", keyGlob).Find(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "list kv")
	}

	return kvs, nil
}

func (r *Repo) Delete(ctx context.Context, appID uuid.UUID, key string) error {
	_, err := gorm.G[ApplicationKV](r.db).
		Where("application_id = ?", appID).
		Where("key like ?", key).
		Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "delete kv")
	}

	return nil
}
