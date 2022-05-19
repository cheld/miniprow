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

package ranch

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cheld/miniprow/pkg/boskos/common"
	"github.com/cheld/miniprow/pkg/boskos/storage"
)

var (
	startTime = fakeTime(time.Now().Truncate(time.Second))
	fakeNow   = fakeTime(startTime.Add(time.Second).Truncate(time.Second))
)

type nameGenerator struct {
	lock  sync.Mutex
	index int
}

func (g *nameGenerator) name() string {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.index++
	return fmt.Sprintf("new-dynamic-res-%d", g.index)
}

func MakeTestRanch(resources []common.Resource) *Ranch {
	persistence := storage.NewMemoryStorage()
	s := NewStorage(persistence)
	s.now = func() time.Time {
		return fakeNow
	}
	nameGen := &nameGenerator{}
	s.generateName = nameGen.name
	ranch := &Ranch{
		Storage: s,
	}
	ranch.now = func() time.Time {
		return fakeNow
	}
	ranch.requestMgr = NewRequestManager(testTTL)
	for _, res := range resources {
		persistence.Add(res, common.NewTenant())
	}

	return ranch
}

// json does not serialized time with nanosecond precision
func fakeTime(t time.Time) time.Time {
	format := "2006-01-02 15:04:05.000"
	now, _ := time.Parse(format, t.Format(format))
	return now
}

func AreErrorsEqual(got error, expect error) bool {
	if got == nil && expect == nil {
		return true
	}

	if got == nil || expect == nil {
		return false
	}

	switch got.(type) {
	default:
		return false
	case *OwnerNotMatch:
		if o, ok := expect.(*OwnerNotMatch); ok {
			if o.request == got.(*OwnerNotMatch).request && o.owner == got.(*OwnerNotMatch).owner {
				return true
			}
		}
		return false
	case *ResourceNotFound:
		if o, ok := expect.(*ResourceNotFound); ok {
			if o.name == got.(*ResourceNotFound).name {
				return true
			}
		}
		return false
	case *ResourceTypeNotFound:
		if o, ok := expect.(*ResourceTypeNotFound); ok {
			if o.rType == got.(*ResourceTypeNotFound).rType {
				return true
			}
		}
		return false
	case *StateNotMatch:
		if o, ok := expect.(*StateNotMatch); ok {
			if o.expect == got.(*StateNotMatch).expect && o.current == got.(*StateNotMatch).current {
				return true
			}
		}
		return false
	}
}

func TestSimple(t *testing.T) {
	c := MakeTestRanch([]common.Resource{})
	c.Acquire("tc.rtype", "tc.state", "tc.dest", "tc.owner", "requestID", common.NewTenant())
}

func TestAcquire(t *testing.T) {
	FakeNow := startTime
	var testcases = []struct {
		name      string
		resources []common.Resource
		owner     string
		rtype     string
		state     string
		dest      string
		expectErr error
	}{
		{
			name:      "ranch has no resource",
			resources: []common.Resource{},
			owner:     "user",
			rtype:     "t",
			state:     "s",
			dest:      "d",
			expectErr: &ResourceTypeNotFound{"t"},
		},
		{
			name: "no match type",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "wrong",
					State:      "s",
					Owner:      "",
					LastUpdate: FakeNow,
				},
			},
			owner:     "user",
			rtype:     "t",
			state:     "s",
			dest:      "d",
			expectErr: &ResourceTypeNotFound{"t"},
		},
		{
			name: "no match state",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "wrong",
					Owner:      "",
					LastUpdate: FakeNow,
				},
			},
			owner:     "user",
			rtype:     "t",
			state:     "s",
			dest:      "d",
			expectErr: &ResourceNotFound{"t"},
		},
		{
			name: "busy",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "foo",
					LastUpdate: FakeNow,
				},
			},
			owner:     "user",
			rtype:     "t",
			state:     "s",
			dest:      "d",
			expectErr: &ResourceNotFound{"t"},
		},
		{
			name: "ok",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "",
					LastUpdate: FakeNow,
				},
			},
			owner:     "user",
			rtype:     "t",
			state:     "s",
			dest:      "d",
			expectErr: nil,
		},
	}

	for _, tc := range testcases {
		c := MakeTestRanch(tc.resources)
		res, _, err := c.Acquire(tc.rtype, tc.state, tc.dest, tc.owner, "requestID", common.NewTenant())
		if !AreErrorsEqual(err, tc.expectErr) {
			t.Errorf("%s - Got error %v, expect error %v", tc.name, err, tc.expectErr)
			continue
		}

		if err == nil {
			if res.State != tc.dest {
				t.Errorf("%s - Wrong final state. Got %v, expect %v", tc.name, res.State, tc.dest)
			}
			resources, _ := c.Storage.GetResources(common.NewTenant())
			if *res != *resources[0] {
				t.Errorf("%s - Wrong resource. Got %v, expect %v", tc.name, res, resources[0])
			} else if !res.LastUpdate.After(FakeNow) {
				t.Errorf("%s - LastUpdate did not update.", tc.name)
			}
		} else {
			resources, _ := c.Storage.GetResources(common.NewTenant())
			for _, res := range resources {
				if res.LastUpdate != FakeNow {
					t.Errorf("%s - LastUpdate should not update. Got %v, expect %v", tc.name, resources[0].LastUpdate, FakeNow)
				}
			}
		}
	}
}

func TestAcquireRoundRobin(t *testing.T) {
	var resources []common.Resource
	for i := 1; i < 5; i++ {
		resources = append(resources, newResource(fmt.Sprintf("res-%d", i), "t", "s", "", startTime))
	}

	results := map[string]int{}

	c := MakeTestRanch(resources)
	for i := 0; i < 4; i++ {
		res, _, err := c.Acquire("t", "s", "d", "foo", "", common.NewTenant())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		_, found := results[res.Name]
		if found {
			t.Errorf("resource %s was used more than once", res.Name)
		}
		c.Release(res.Name, "s", "foo", common.NewTenant())
	}
}

func newResource(name, t, state, owner string, lastUpdate time.Time) common.Resource {
	return common.Resource{
		Name:       name,
		Type:       t,
		State:      state,
		Owner:      owner,
		LastUpdate: lastUpdate,
	}
}

func TestAcquireOrder(t *testing.T) {
	FakeNow := startTime
	resources := []common.Resource{
		{
			Name:       "res-1",
			Type:       "t",
			State:      "s",
			Owner:      "",
			LastUpdate: FakeNow,
		},
		{
			Name:       "res-2",
			Type:       "t",
			State:      "s",
			Owner:      "",
			LastUpdate: FakeNow,
		},
		{
			Name:       "res-3",
			Type:       "t",
			State:      "s",
			Owner:      "",
			LastUpdate: FakeNow,
		},
		{
			Name:       "res-4",
			Type:       "t",
			State:      "s",
			Owner:      "",
			LastUpdate: FakeNow,
		},
	}

	expected := []string{"res-1", "res-2", "res-3", "res-4"}

	c := MakeTestRanch(resources)
	for i := 0; i < 2; i++ {
		_, _, err := c.Acquire("t", "s", "d", "foo", "requestID", common.NewTenant())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	updatedRessources, _ := c.Storage.GetResources(common.NewTenant())
	for idx, _ := range updatedRessources {
		if updatedRessources[idx].Name != expected[idx] {
			t.Errorf("Resource %d, expected %v, got %v", idx, expected[idx], updatedRessources[idx].Name)
		}
	}

}

func TestAcquireOnDemand(t *testing.T) {
	owner := "tester"
	rType := "dr"
	requestID1 := "req1234"
	requestID2 := "req12345"
	requestID3 := "req123456"
	duration := time.Second
	config := &common.BoskosConfig{
		Resources: []common.ResourceEntry{
			{
				Type:     rType,
				State:    common.Dirty,
				MinCount: 0,
				MaxCount: 2,
				LifeSpan: &common.Duration{&duration},
			},
		},
	}

	// Not adding any resources to start with
	c := MakeTestRanch([]common.Resource{})
	c.Storage.SyncResources(config, common.NewTenant())

	// First acquire should trigger a creation
	if _, _, err := c.Acquire(rType, common.Free, common.Busy, owner, requestID1, common.NewTenant()); err == nil {
		t.Errorf("should fail since there is not resource yet")
	}
	if resources, err := c.Storage.GetResources(common.NewTenant()); err != nil {
		t.Fatal(err)
	} else if len(resources) != 1 {
		t.Fatal("A resource should have been created")
	}
	// Attempting to create another resource
	if _, _, err := c.Acquire(rType, common.Free, common.Busy, owner, requestID1, common.NewTenant()); err == nil {
		t.Errorf("should succeed since the created is dirty")
	}
	if resources, err := c.Storage.GetResources(common.NewTenant()); err != nil {
		t.Error(err)
	} else if len(resources) != 1 {
		t.Errorf("No new resource should have been created")
	}
	// Creating another
	if _, _, err := c.Acquire(rType, common.Free, common.Busy, owner, requestID2, common.NewTenant()); err == nil {
		t.Errorf("should succeed since the created is dirty")
	}
	if resources, err := c.Storage.GetResources(common.NewTenant()); err != nil {
		t.Error(err)
	} else if len(resources) != 2 {
		t.Errorf("Another resource should have been created")
	}
	// Attempting to create another
	if _, _, err := c.Acquire(rType, common.Free, common.Busy, owner, requestID3, common.NewTenant()); err == nil {
		t.Errorf("should fail since there is not resource yet")
	}
	resources, err := c.Storage.GetResources(common.NewTenant())
	if err != nil {
		t.Error(err)
	} else if len(resources) != 2 {
		t.Errorf("No other resource should have been created")
	}
	for _, res := range resources {
		c.Storage.DeleteResource(res.Name, common.NewTenant())
	}
	if _, _, err := c.Acquire(rType, common.Free, common.Busy, owner, "", common.NewTenant()); err == nil {
		t.Errorf("should fail since there is not resource yet")
	}
	if resources, err := c.Storage.GetResources(common.NewTenant()); err != nil {
		t.Error(err)
	} else if len(resources) != 0 {
		t.Errorf("No new resource should have been created")
	}
}

func TestRelease(t *testing.T) {
	FakeNow := time.Now()
	var testcases = []struct {
		name      string
		resources []common.Resource
		resName   string
		owner     string
		dest      string
		expectErr error
	}{
		{
			name:      "ranch has no resource",
			resources: []common.Resource{},
			resName:   "res",
			owner:     "user",
			dest:      "d",
			expectErr: &ResourceNotFound{"res"},
		},
		{
			name: "wrong owner",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName:   "res",
			owner:     "user",
			dest:      "d",
			expectErr: &OwnerNotMatch{"user", "merlin"},
		},
		{
			name: "no match name",
			resources: []common.Resource{
				{
					Name:       "foo",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName:   "res",
			owner:     "user",
			dest:      "d",
			expectErr: &ResourceNotFound{"res"},
		},
		{
			name: "ok",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName:   "res",
			owner:     "merlin",
			dest:      "d",
			expectErr: nil,
		},
	}

	for _, tc := range testcases {
		c := MakeTestRanch(tc.resources)
		err := c.Release(tc.resName, tc.dest, tc.owner, common.NewTenant())
		if !AreErrorsEqual(err, tc.expectErr) {
			t.Errorf("%s - Got error %v, expect error %v", tc.name, err, tc.expectErr)
			continue
		}

		if err == nil {
			resources, _ := c.Storage.GetResources(common.NewTenant())
			if resources[0].Owner != "" {
				t.Errorf("%s - Wrong owner after release. Got %v, expect empty", tc.name, resources[0].Owner)
			} else if resources[0].State != tc.dest {
				t.Errorf("%s - Wrong state after release. Got %v, expect %v", tc.name, resources[0].State, tc.dest)
			} else if !resources[0].LastUpdate.After(FakeNow) {
				t.Errorf("%s - LastUpdate did not update.", tc.name)
			}
		} else {
			resources, _ := c.Storage.GetResources(common.NewTenant())
			for _, res := range resources {
				if res.LastUpdate != FakeNow {
					t.Errorf("%s - LastUpdate should not update. Got %v, expect %v", tc.name, resources[0].LastUpdate, FakeNow)
				}
			}
		}
	}
}

func TestReset(t *testing.T) {
	FakeNow := startTime

	var testcases = []struct {
		name       string
		resources  []common.Resource
		rtype      string
		state      string
		dest       string
		expire     time.Duration
		hasContent bool
	}{

		{
			name: "empty - has no owner",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "",
					LastUpdate: FakeNow.Add(-time.Minute * 20),
				},
			},
			rtype:  "t",
			state:  "s",
			expire: time.Minute * 10,
			dest:   "d",
		},
		{
			name: "empty - not expire",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "",
					LastUpdate: FakeNow,
				},
			},
			rtype:  "t",
			state:  "s",
			expire: time.Minute * 10,
			dest:   "d",
		},
		{
			name: "empty - no match type",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "wrong",
					State:      "s",
					Owner:      "",
					LastUpdate: FakeNow.Add(-time.Minute * 20),
				},
			},
			rtype:  "t",
			state:  "s",
			expire: time.Minute * 10,
			dest:   "d",
		},
		{
			name: "empty - no match state",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "wrong",
					Owner:      "",
					LastUpdate: FakeNow.Add(-time.Minute * 20),
				},
			},
			rtype:  "t",
			state:  "s",
			expire: time.Minute * 10,
			dest:   "d",
		},
		{
			name: "ok",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "user",
					LastUpdate: FakeNow.Add(-time.Minute * 20),
				},
			},
			rtype:      "t",
			state:      "s",
			expire:     time.Minute * 10,
			dest:       "d",
			hasContent: true,
		},
	}

	for _, tc := range testcases {
		c := MakeTestRanch(tc.resources)
		rmap, _ := c.Reset(tc.rtype, tc.state, tc.expire, tc.dest, common.NewTenant())

		if !tc.hasContent {
			if len(rmap) != 0 {
				t.Errorf("%s - Expect empty map. Got %v", tc.name, rmap)
			}
		} else {
			if owner, ok := rmap["res"]; !ok || owner != "user" {
				t.Errorf("%s - Expect res - user. Got %v", tc.name, rmap)
			}
			resources, _ := c.Storage.GetResources(common.NewTenant())
			fmt.Println(FakeNow)
			fmt.Println(fakeTime(fakeNow))
			fmt.Println(resources[0].LastUpdate)
			if !resources[0].LastUpdate.After(FakeNow) {
				t.Errorf("%s - LastUpdate did not update.", tc.name)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	FakeNow := startTime

	var testcases = []struct {
		name      string
		resources []common.Resource
		resName   string
		owner     string
		state     string
		expectErr error
	}{
		{
			name:      "ranch has no resource",
			resources: []common.Resource{},
			resName:   "res",
			owner:     "user",
			state:     "s",
			expectErr: &ResourceNotFound{"res"},
		},
		{
			name: "wrong owner",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName:   "res",
			owner:     "user",
			state:     "s",
			expectErr: &OwnerNotMatch{"user", "merlin"},
		},
		{
			name: "wrong state",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName:   "res",
			owner:     "merlin",
			state:     "foo",
			expectErr: &StateNotMatch{"s", "foo"},
		},
		{
			name: "no matched resource",
			resources: []common.Resource{
				{
					Name:       "foo",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName:   "res",
			owner:     "merlin",
			state:     "s",
			expectErr: &ResourceNotFound{"res"},
		},
		{
			name: "ok",
			resources: []common.Resource{
				{
					Name:       "res",
					Type:       "t",
					State:      "s",
					Owner:      "merlin",
					LastUpdate: FakeNow,
				},
			},
			resName: "res",
			owner:   "merlin",
			state:   "s",
		},
	}

	for _, tc := range testcases {
		c := MakeTestRanch(tc.resources)
		err := c.Update(tc.resName, tc.owner, tc.state, nil, common.NewTenant())
		if !AreErrorsEqual(err, tc.expectErr) {
			t.Errorf("%s - Got error %v, expect error %v", tc.name, err, tc.expectErr)
			continue
		}

		if err == nil {
			resources, _ := c.Storage.GetResources(common.NewTenant())
			if resources[0].Owner != tc.owner {
				t.Errorf("%s - Wrong owner after release. Got %v, expect %v", tc.name, resources[0].Owner, tc.owner)
			} else if resources[0].State != tc.state {
				t.Errorf("%s - Wrong state after release. Got %v, expect %v", tc.name, resources[0].State, tc.state)
			} else if !resources[0].LastUpdate.After(FakeNow) {
				t.Errorf("%s - LastUpdate did not update.", tc.name)
			}
		} else {
			resources, _ := c.Storage.GetResources(common.NewTenant())
			for _, res := range resources {
				if res.LastUpdate != FakeNow {
					t.Errorf("%s - LastUpdate should not update. Got %v, expect %v", tc.name, resources[0].LastUpdate, FakeNow)
				}
			}
		}
	}
}

func TestMetric(t *testing.T) {
	var testcases = []struct {
		name         string
		resources    []common.Resource
		metricType   string
		expectErr    error
		expectMetric common.Metric
	}{
		{
			name:       "ranch has no resource",
			resources:  []common.Resource{},
			metricType: "t",
			expectErr:  &ResourceNotFound{"t"},
		},
		{
			name: "no matching resource",
			resources: []common.Resource{
				{
					Name:  "res",
					Type:  "t",
					State: "s",
					Owner: "merlin",
				},
			},
			metricType: "foo",
			expectErr:  &ResourceNotFound{"foo"},
		},
		{
			name: "one resource",
			resources: []common.Resource{
				{
					Name:  "res",
					Type:  "t",
					State: "s",
					Owner: "merlin",
				},
			},
			metricType: "t",
			expectMetric: common.Metric{
				Type: "t",
				Current: map[string]int{
					"s": 1,
				},
				Owners: map[string]int{
					"merlin": 1,
				},
			},
		},
		{
			name: "multiple resources",
			resources: []common.Resource{
				{
					Name:  "res-1",
					Type:  "t",
					State: "s",
					Owner: "merlin",
				},
				{
					Name:  "res-2",
					Type:  "t",
					State: "p",
					Owner: "pony",
				},
				{
					Name:  "res-2",
					Type:  "t",
					State: "s",
					Owner: "pony",
				},
				{
					Name:  "res-3",
					Type:  "foo",
					State: "s",
					Owner: "pony",
				},
				{
					Name:  "res-4",
					Type:  "t",
					State: "d",
					Owner: "merlin",
				},
			},
			metricType: "t",
			expectMetric: common.Metric{
				Type: "t",
				Current: map[string]int{
					"s": 1,
					"d": 1,
					"p": 1,
				},
				Owners: map[string]int{
					"merlin": 2,
					"pony":   1,
				},
			},
		},
	}

	for _, tc := range testcases {
		c := MakeTestRanch(tc.resources)
		metric, err := c.Metric(tc.metricType, common.NewTenant())
		if !AreErrorsEqual(err, tc.expectErr) {
			t.Errorf("%s - Got error %v, expect error %v", tc.name, err, tc.expectErr)
			//continue
		}

		if err == nil {
			if !reflect.DeepEqual(metric, tc.expectMetric) {
				t.Errorf("%s - wrong metric, got %v, want %v", tc.name, metric, tc.expectMetric)
			}
		}
	}
}

/*
func TestSyncConfig(t *testing.T) {
	var testcases = []struct {
		name   string
		oldRes []common.Resource
		newRes []common.Resource
		expect []common.Resource
	}{
		{
			name:   "empty",
			oldRes: []common.Resource{},
			newRes: []common.Resource{},
			expect: []common.Resource{},
		},
		{
			name:   "append",
			oldRes: []common.Resource{},
			newRes: []common.Resource{
				{
					Name: "res",
					Type: "t",
				},
			},
			expect: []common.Resource{
				{
					Name:  "res",
					Type:  "t",
					State: "free",
				},
			},
		},
		{
			name: "should not have a type change",
			oldRes: []common.Resource{
				{
					Name: "res",
					Type: "t",
				},
			},
			newRes: []common.Resource{
				{
					Name: "res",
					Type: "d",
				},
			},
			expect: []common.Resource{
				{
					Name: "res",
					Type: "t",
				},
			},
		},
		{
			name: "delete",
			oldRes: []common.Resource{
				{
					Name: "res",
					Type: "t",
				},
			},
			newRes: []common.Resource{},
			expect: []common.Resource{},
		},
		{
			name: "delete busy",
			oldRes: []common.Resource{
				{
					Name:  "res",
					Type:  "t",
					State: "busy",
					Owner: "o",
				},
			},
			newRes: []common.Resource{},
			expect: []common.Resource{
				{
					Name:  "res",
					Type:  "t",
					State: "busy",
					Owner: "o",
				},
			},
		},
		{
			name: "append and delete",
			oldRes: []common.Resource{
				{
					Name: "res-1",
					Type: "t",
				},
			},
			newRes: []common.Resource{
				{
					Name: "res-2",
					Type: "t",
				},
			},
			expect: []common.Resource{
				{
					Name:  "res-2",
					Type:  "t",
					State: "free",
				},
			},
		},
		{
			name: "append and delete busy",
			oldRes: []common.Resource{
				{
					Name:  "res-1",
					Type:  "t",
					State: "busy",
					Owner: "o",
				},
			},
			newRes: []common.Resource{
				{
					Name: "res-2",
					Type: "t",
				},
			},
			expect: []common.Resource{
				{
					Name:  "res-1",
					Type:  "t",
					State: "busy",
					Owner: "o",
				},
				{
					Name:  "res-2",
					Type:  "t",
					State: "free",
				},
			},
		},
		{
			name: "append/delete mixed type",
			oldRes: []common.Resource{
				{
					Name: "res-1",
					Type: "t",
				},
			},
			newRes: []common.Resource{
				{
					Name: "res-2",
					Type: "t",
				},
				{
					Name: "res-3",
					Type: "t2",
				},
			},
			expect: []common.Resource{
				{
					Name:  "res-2",
					Type:  "t",
					State: "free",
				},
				{
					Name:  "res-3",
					Type:  "t2",
					State: "free",
				},
			},
		},
	}

	for _, tc := range testcases {
		c := MakeTestRanch(tc.oldRes)
		//c.syncConfigHelper(tc.newRes)
		resources, _ := c.Storage.GetResources()
		if !reflect.DeepEqual(resources, tc.expect) {
			t.Errorf("Test %v: got %v, expect %v", tc.name, resources, tc.expect)
		}
	}
}
*/
