package managerservice

import (
	"context"
	"vibesiter/pkg/dto"

	"github.com/cockroachdb/errors"
)

func (r *Service) HandleHTTP(ctx context.Context, appSlug string, req *dto.HTTPRequest) (dto.HTTPResponse, error) {
	app, err := r.managerRepo.SelectApplicationBySlug(ctx, appSlug)
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "select application by slug")
	}

	resp, err := r.llmSupplier.HTTP(ctx, &app, req)
	if err != nil {
		return dto.HTTPResponse{}, errors.Wrap(err, "http request")
	}

	return resp, nil
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
