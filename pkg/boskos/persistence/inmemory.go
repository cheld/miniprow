package persistence

import (
	"fmt"
	"sort"
	"sync"

	"github.com/cheld/miniprow/pkg/boskos/common"
)

type inMemoryStore struct {
	resources map[string]map[string]common.Resource
	lock      sync.RWMutex
}

// NewResourceMemoryStorage creates an in memory persistence layer
func NewResourceMemoryStorage() ResourcePersistence {
	mem := map[string]map[string]common.Resource{}
	mem["default"] = map[string]common.Resource{}
	return &inMemoryStore{
		resources: mem,
	}
}

func (im *inMemoryStore) Add(r common.Resource, tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	if im.resources[tenant.ID()] == nil {
		im.resources[tenant.ID()] = map[string]common.Resource{}
	}
	_, ok := im.resources[tenant.ID()][r.Name]
	if ok {
		return fmt.Errorf("resource %s already exists", r.Name)
	}
	im.resources[tenant.ID()][r.Name] = r
	return nil
}

func (im *inMemoryStore) Delete(name string, tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	if im.resources[tenant.ID()] == nil {
		return fmt.Errorf("cannot find item %s", name)
	}
	_, ok := im.resources[tenant.ID()][name]
	if !ok {
		return fmt.Errorf("cannot find item %s", name)
	}
	delete(im.resources[tenant.ID()], name)
	return nil
}

func (im *inMemoryStore) Update(r common.Resource, tenant common.Tenant) (common.Resource, error) {
	im.lock.Lock()
	defer im.lock.Unlock()
	if im.resources[tenant.ID()] == nil {
		return common.Resource{}, fmt.Errorf("cannot find item %s", r.Name)
	}
	_, ok := im.resources[tenant.ID()][r.Name]
	if !ok {
		return common.Resource{}, fmt.Errorf("cannot find item %s", r.Name)
	}
	im.resources[tenant.ID()][r.Name] = r
	return r, nil
}

func (im *inMemoryStore) Get(name string, tenant common.Tenant) (common.Resource, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	if im.resources[tenant.ID()] == nil {
		return common.Resource{}, fmt.Errorf("cannot find item %s", name)
	}
	r, ok := im.resources[tenant.ID()][name]
	if !ok {
		return common.Resource{}, fmt.Errorf("cannot find item %s", name)
	}
	return r, nil
}

func (im *inMemoryStore) List(tenant common.Tenant) ([]common.Resource, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	if im.resources[tenant.ID()] == nil {
		return []common.Resource{}, nil
	}
	var resources []common.Resource
	for _, r := range im.resources[tenant.ID()] {
		resources = append(resources, r)
	}
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})
	return resources, nil
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
