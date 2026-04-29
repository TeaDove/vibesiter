package managerservice

import (
	"vibesiter/pkg/kvrepo"
	"vibesiter/pkg/llmsupplier"
	"vibesiter/pkg/managerrepo"
)

type Service struct {
	managerRepo *managerrepo.Repo
	kvRepo      *kvrepo.Repo
	llmSupplier *llmsupplier.Supplier
}

func NewService(managerRepo *managerrepo.Repo, kvRepo *kvrepo.Repo, llmSupplier *llmsupplier.Supplier) *Service {
	return &Service{managerRepo, kvRepo, llmSupplier}
}
