package persistence

import (
	"sync"

	"github.com/cheld/miniprow/pkg/boskos/common"
)

type TenantPersistence interface {
	AddToken(token string, tenant common.Tenant) error
	DeleteToken(tenant common.Tenant) error
	GetTenant(token, project string) (common.Tenant, error)
}

type tenantInMemoryStore struct {
	tokens map[string]string
	lock   sync.RWMutex
}

// NewMemoryStorage creates an in memory persistence layer
func NewTenantMemoryStorage() TenantPersistence {
	mem := map[string]string{}
	return &tenantInMemoryStore{
		tokens: mem,
	}
}

func (im *tenantInMemoryStore) AddToken(token string, tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	return nil
}

func (im *tenantInMemoryStore) DeleteToken(tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	return nil
}

func (im *tenantInMemoryStore) GetTenant(token, project string) (common.Tenant, error) {
	im.lock.Lock()
	defer im.lock.Unlock()
	return common.NewTenant(), nil
}
