package managerservice

import (
	"vibesiter/pkg/llmsupplier"
	"vibesiter/pkg/managerrepo"
)

type Service struct {
	managerRepo *managerrepo.Repo
	llmSupplier *llmsupplier.Supplier
}

func NewService(managerRepo *managerrepo.Repo, llmSupplier *llmsupplier.Supplier) *Service {
	return &Service{managerRepo, llmSupplier}
}
