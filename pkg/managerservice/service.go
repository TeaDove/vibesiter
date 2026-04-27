package managerservice

import (
	"context"
	"ws-lan-chat/pkg/managerrepo"

	"github.com/cockroachdb/errors"
)

type Service struct {
	managerRepo *managerrepo.Repo
}

func NewService(managerRepo *managerrepo.Repo) *Service {
	return &Service{managerRepo}
}

func (r *Service) GenerateApp(_ context.Context, userPrompt string) (managerrepo.Application, error) {
	if userPrompt == "" {
		return managerrepo.Application{}, errors.New("userPrompt is required")
	}

	return managerrepo.Application{}, nil
}
