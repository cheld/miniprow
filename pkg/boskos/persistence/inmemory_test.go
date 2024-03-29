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

// This is an xtest because it imports the crds package, but the crds package
// also imports storage, creating a cycle.
package persistence

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/cheld/miniprow/pkg/boskos/common"
	"github.com/cheld/miniprow/pkg/common/core"
)

func createStorages() []Persistence {
	return []Persistence{
		NewPersistence(),
	}
}

func TestAddDelete(t *testing.T) {
	for _, s := range createStorages() {
		var resources []common.Resource
		var err error
		for i := 0; i < 10; i++ {
			resources = append(resources, common.Resource{
				Name: fmt.Sprintf("res-%d", i),
				Type: fmt.Sprintf("type-%d", i),
			})
		}
		sort.Stable(common.ResourceByName(resources))
		for _, res := range resources {
			if err = s.Add(res, core.NewTenant()); err != nil {
				t.Errorf("unable to add %s, %v", res.Name, err)
			}
		}
		returnedResources, err := s.List(core.NewTenant())
		if err != nil {
			t.Errorf("unable to list resources, %v", err)
		}
		sort.Stable(common.ResourceByName(returnedResources))
		if !reflect.DeepEqual(resources, returnedResources) {
			t.Errorf("received resources (%v) do not match resources (%v)", resources, returnedResources)
		}
		for _, r := range returnedResources {
			err = s.Delete(r.Name, core.NewTenant())
			if err != nil {
				t.Errorf("unable to delete resource %s.%v", r.Name, err)
			}
		}
		eResources, err := s.List(core.NewTenant())
		if err != nil {
			t.Errorf("unable to list resources, %v", err)
		}
		if len(eResources) != 0 {
			t.Error("list should return an empty list")
		}
	}
}

func TestUpdateGet(t *testing.T) {
	for _, s := range createStorages() {
		oRes := common.Resource{
			Name: "original",
			Type: "type",
		}
		if err := s.Add(oRes, core.NewTenant()); err != nil {
			t.Errorf("unable to add resource, %v", err)
		}
		uRes := oRes
		uRes.Type = "typeUpdated"
		if _, err := s.Update(uRes, core.NewTenant()); err != nil {
			t.Errorf("unable to update resource %v", err)
		}
		res, err := s.Get(oRes.Name, core.NewTenant())
		if err != nil {
			t.Errorf("unable to get resource, %v", err)
		}
		if !reflect.DeepEqual(uRes, res) {
			t.Errorf("expected (%v) and received (%v) do not match", uRes, res)
		}
	}
}

func TestNegativeDeleteGet(t *testing.T) {
	for _, s := range createStorages() {
		oRes := common.Resource{
			Name: "original",
			Type: "type",
		}
		if err := s.Add(oRes, core.NewTenant()); err != nil {
			t.Errorf("unable to add resource, %v", err)
		}
		uRes := common.Resource{
			Name: "notExist",
			Type: "type",
		}
		if _, err := s.Update(uRes, core.NewTenant()); err == nil {
			t.Errorf("should not be able to update resource, %v", err)
		}
		if err := s.Delete(uRes.Name, core.NewTenant()); err == nil {
			t.Errorf("should not be able to delete resource, %v", err)
		}
	}
}

func TestMultiTenant(t *testing.T) {
	for _, s := range createStorages() {
		oRes := common.Resource{
			Name: "original",
			Type: "type",
		}
		if err := s.Add(oRes, core.NewTenant()); err != nil {
			t.Errorf("unable to add resource, %v", err)
		}
		_, err := s.Get(oRes.Name, core.NewTenant())
		if err != nil {
			t.Errorf("unable to get resource, %v", err)
		}
		_, err = s.Get(oRes.Name, core.Tenant{
			Organization: "different",
			Project:      "different",
		})
		if err == nil {
			t.Errorf("should not be able to retrieve resource from different tenant, %v", err)
		}
	}
}
