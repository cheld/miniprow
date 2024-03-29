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

package persistence

import (
	"github.com/cheld/miniprow/pkg/boskos/common"
	"github.com/cheld/miniprow/pkg/common/core"
)

// PersistenceLayer defines a simple interface to persists Boskos Information
type ClientCache interface {
	Add(r common.Resource, tenant core.Tenant) error
	Delete(name string, tenant core.Tenant) error
	Update(r common.Resource, tenant core.Tenant) (common.Resource, error)
	Get(name string, tenant core.Tenant) (common.Resource, error)
	List(tenant core.Tenant) ([]common.Resource, error)
}

type Persistence interface {
	Add(r common.Resource, tenant core.Tenant) error
	Delete(name string, tenant core.Tenant) error
	Update(r common.Resource, tenant core.Tenant) (common.Resource, error)
	Get(name string, tenant core.Tenant) (common.Resource, error)
	List(tenant core.Tenant) ([]common.Resource, error)
	AddDynamicResourceLifeCycle(r common.DynamicResourceLifeCycle, tenant core.Tenant) error
	GetDynamicResourceLifeCycle(rtype string, tenant core.Tenant) (common.DynamicResourceLifeCycle, error)
	AddToken(token string, tenant core.Tenant) error
	DeleteToken(tenant core.Tenant) error
	GetTenantFromToken(token, project string) (core.Tenant, error)
	Close()
}
