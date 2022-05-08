/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"sync"

	"fmt"

	"github.com/cheld/miniprow/pkg/boskos/common"
)

// PersistenceLayer defines a simple interface to persists Boskos Information
type PersistenceLayer interface {
	Add(r common.Resource, organization, project string) error
	Delete(name string, organization, project string) error
	Update(r common.Resource, organization, project string) (common.Resource, error)
	Get(name string, organization, project string) (common.Resource, error)
	List(organization, project string) ([]common.Resource, error)
}

type inMemoryStore struct {
	resources map[string]map[string]common.Resource
	lock      sync.RWMutex
}

// NewMemoryStorage creates an in memory persistence layer
func NewMemoryStorage() PersistenceLayer {
	mem := map[string]map[string]common.Resource{}
	mem["default"] = map[string]common.Resource{}
	return &inMemoryStore{
		resources: mem,
	}
}

func (im *inMemoryStore) Add(r common.Resource, org, project string) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	tenant := org + project
	if im.resources[tenant] == nil {
		im.resources[tenant] = map[string]common.Resource{}
	}
	_, ok := im.resources[tenant][r.Name]
	if ok {
		return fmt.Errorf("resource %s already exists", r.Name)
	}
	im.resources[tenant][r.Name] = r
	return nil
}

func (im *inMemoryStore) Delete(name string, org, project string) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	tenant := org + project
	if im.resources[tenant] == nil {
		return fmt.Errorf("cannot find item %s", name)
	}
	_, ok := im.resources[tenant][name]
	if !ok {
		return fmt.Errorf("cannot find item %s", name)
	}
	delete(im.resources[tenant], name)
	return nil
}

func (im *inMemoryStore) Update(r common.Resource, org, project string) (common.Resource, error) {
	im.lock.Lock()
	defer im.lock.Unlock()
	tenant := org + project
	if im.resources[tenant] == nil {
		return common.Resource{}, fmt.Errorf("cannot find item %s", r.Name)
	}
	_, ok := im.resources[tenant][r.Name]
	if !ok {
		return common.Resource{}, fmt.Errorf("cannot find item %s", r.Name)
	}
	im.resources[tenant][r.Name] = r
	return r, nil
}

func (im *inMemoryStore) Get(name, org, project string) (common.Resource, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	tenant := org + project
	if im.resources[tenant] == nil {
		return common.Resource{}, fmt.Errorf("cannot find item %s", name)
	}
	r, ok := im.resources[tenant][name]
	if !ok {
		return common.Resource{}, fmt.Errorf("cannot find item %s", name)
	}
	return r, nil
}

func (im *inMemoryStore) List(org, project string) ([]common.Resource, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	tenant := org + project
	if im.resources[tenant] == nil {
		return []common.Resource{}, nil
	}
	var resources []common.Resource
	for _, r := range im.resources[tenant] {
		resources = append(resources, r)
	}
	return resources, nil
}
