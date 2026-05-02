package managerservice

import (
	"context"
	"vibesiter/pkg/dto"
	"vibesiter/pkg/kvrepo"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (r *Service) HandleHTTP(ctx context.Context, appSlug string, req *dto.HTTPRequest) (dto.HTTPResponse, error) {
	app, err := r.managerRepo.SelectApplicationBySlug(ctx, appSlug)
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "select application by slug")
	}

	zerolog.Ctx(ctx).Debug().
		Str("slug", appSlug).
		Interface("req", req).
		Msg("handling.llm.http.request")

	resp, messages, err := r.llmSupplier.HTTP(ctx, &app, req)
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "http request")
	}

	var kvs []kvrepo.ApplicationKV

	for {
		kvs, err = r.executeActions(ctx, app.ID, resp.Actions)
		if err != nil {
			return dto.HTTPResponse{}, errors.Wrap(err, "execute actions")
		}

		if resp.Response != nil {
			return *resp.Response, nil
		}

		resp, messages, err = r.llmSupplier.HTTPContinue(ctx, messages, kvs)
		if err != nil {
			return dto.HTTPResponse{}, errors.Wrap(err, "http continue")
		}
	}
}

func (r *Service) executeActions(
	ctx context.Context,
	appID uuid.UUID,
	actions []dto.ActionKV,
) ([]kvrepo.ApplicationKV, error) {
	if len(actions) == 0 {
		return nil, nil
	}

	var keys []kvrepo.ApplicationKV

	for _, action := range actions {
		switch action.Type {
		case dto.ActionTypeKvSet:
			err := r.kvRepo.Set(ctx, appID, action.Key, action.Value)
			if err != nil {
				return nil, errors.Wrap(err, "set kv")
			}
		case dto.ActionTypeKvDelete:
			err := r.kvRepo.Delete(ctx, appID, action.Key)
			if err != nil {
				return nil, errors.Wrap(err, "set kv")
			}
		case dto.ActionTypeKvList:
			kvs, err := r.kvRepo.List(ctx, appID, action.Key)
			if err != nil {
				return nil, errors.Wrap(err, "list kv")
			}

			keys = append(keys, kvs...)
		case dto.ActionTypeKvGet:
			kv, err := r.kvRepo.Get(ctx, appID, action.Key)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.Wrap(err, "get kv")
			}

			keys = append(keys, kv)
		}
	}

	return keys, nil
}

func (r *Service) Serve(ctx context.Context, appSlug string, path string) ([]byte, error) {
	app, err := r.managerRepo.SelectApplicationBySlug(ctx, appSlug)
	if err != nil {
		return nil, errors.Wrap(err, "select application by slug")
	}

	file, err := r.managerRepo.SelectFile(ctx, app.ID, path)
	if err != nil {
		return nil, errors.Wrap(err, "select file")
	}

	return file.Content, nil
}
