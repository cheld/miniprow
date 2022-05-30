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
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cheld/miniprow/pkg/boskos/common"
	"github.com/cheld/miniprow/pkg/common/util"
)

// Ranch is the place which all of the Resource objects lives.
type Ranch struct {
	Storage    *Storage
	requestMgr *RequestManager
	//
	now func() time.Time
}

// Public errors:

// ResourceNotFound will be returned if requested resource does not exist.
type ResourceNotFound struct {
	name string
}

func (r ResourceNotFound) Error() string {
	return fmt.Sprintf("no available resource %s, try again later.", r.name)
}

// ResourceTypeNotFound will be returned if requested resource type does not exist.
type ResourceTypeNotFound struct {
	rType string
}

func (r ResourceTypeNotFound) Error() string {
	return fmt.Sprintf("resource type %q does not exist", r.rType)
}

// OwnerNotMatch will be returned if request owner does not match current owner for target resource.
type OwnerNotMatch struct {
	request string
	owner   string
}

func (o OwnerNotMatch) Error() string {
	return fmt.Sprintf("owner mismatch request by %s, currently owned by %s", o.request, o.owner)
}

// StateNotMatch will be returned if requested state does not match current state for target resource.
type StateNotMatch struct {
	expect  string
	current string
}

func (s StateNotMatch) Error() string {
	return fmt.Sprintf("state mismatch - expected %v, current %v", s.expect, s.current)
}

// NewRanch creates a new Ranch object.
// In: config - path to resource file
//     storage - path to where to save/restore the state data
// Out: A Ranch object, loaded from config/storage, or error
func NewRanch(config *[]byte, s *Storage, ttl time.Duration, tenant common.Tenant) (*Ranch, error) {
	newRanch := &Ranch{
		Storage:    s,
		requestMgr: NewRequestManager(ttl),
		now:        time.Now,
	}
	if err := newRanch.SyncConfig(config, tenant); err != nil {
		return nil, err
	}
	logrus.Infof("Loaded Boskos configuration successfully")
	return newRanch, nil
}

// acquireRequestPriorityKey is used as key for request priority cache.
type acquireRequestPriorityKey struct {
	rType, state, tenant string
}

// Acquire checks out a type of resource in certain state without an owner,
// and move the checked out resource to the end of the resource list.
// In: rtype - name of the target resource
//     state - current state of the requested resource
//     dest - destination state of the requested resource
//     owner - requester of the resource
//     requestID - request ID to get a priority in the queue
// Out: A valid Resource object and the time when the resource was originally requested on success, or
//      ResourceNotFound error if target type resource does not exist in target state.
func (r *Ranch) Acquire(rType, state, dest, owner, requestID string, tenant common.Tenant) (*common.Resource, time.Time, error) {
	logger := logrus.WithFields(logrus.Fields{
		"type":       rType,
		"state":      state,
		"dest":       dest,
		"owner":      owner,
		"identifier": requestID,
	})

	var returnRes *common.Resource
	createdTime := r.now()
	if err := retryOnConflict(func() error {
		logger.Debug("Determining request priority...")
		ts := acquireRequestPriorityKey{rType: rType, state: state, tenant: tenant.ID()}
		rank, new := r.requestMgr.GetRank(ts, requestID)
		logger.WithFields(logrus.Fields{"rank": rank, "new": new}).Debug("Determined request priority.")

		resources, err := r.Storage.GetResources(tenant)
		if err != nil {
			logger.WithError(err).Errorf("could not get resources")
			return &ResourceNotFound{rType}
		}
		logger.Debugf("Considering %d resources.", len(resources))

		// For request priority we need to go over all the list until a matching rank
		matchingResoucesCount := 0
		typeCount := 0
		for _, res := range resources {
			fmt.Printf("Resource %v\n", res)
		}
		for _, res := range resources {
			fmt.Printf("Typecount, Matchcount, rank %v, %v, %v\n", typeCount, matchingResoucesCount, rank)
			if rType != res.Type {
				continue
			}
			typeCount++

			if state != res.State || res.Owner != "" {
				continue
			}
			matchingResoucesCount++

			if matchingResoucesCount < rank {
				continue
			}
			logger = logger.WithField("resource", res.Name)
			res.Owner = owner
			res.State = dest
			logger.Debug("Updating resource.")
			updatedRes, err := r.Storage.UpdateResource(res, tenant)
			if err != nil {
				return err
			}
			// Deleting this request since it has been fulfilled
			if requestID != "" {
				if createdTime, err = r.requestMgr.GetCreatedAt(ts, requestID); err != nil {
					// It is chosen NOT to fail the function since the resource has been already updated to give ownership.
					logger.WithError(err).Errorf("Error occurred when getting the created time")
				}
				logger.Debug("Cleaning up requests.")
				r.requestMgr.Delete(ts, requestID)
			}
			logger.Debug("Successfully acquired resource.")
			returnRes = updatedRes
			return nil
		}

		addResource(new, logger, r, rType, typeCount, tenant)
		fmt.Printf("Typecount: %v", typeCount)
		if typeCount > 0 {
			return &ResourceNotFound{rType}
		}
		return &ResourceTypeNotFound{rType}
	}); err != nil {
		switch err.(type) {
		case *ResourceNotFound:
			// This error occurs when there are no more resources to lease out.
			// Such a condition is a normal and expected part of operation, so
			// it does not warrant an error log.
		default:
			logrus.WithError(err).Error("Acquire failed")
		}
		return nil, createdTime, err
	}

	return returnRes, createdTime, nil
}

func addResource(new bool, logger *logrus.Entry, r *Ranch, rType string, typeCount int, tenant common.Tenant) {
	fmt.Printf("add resource %v\n", new)
	if !new {
		return
	}
	logger.Debug("Checking for associated dynamic resource type...")
	lifeCycle, err := r.Storage.GetDynamicResourceLifeCycle(rType, tenant)
	// // Assuming error means no associated dynamic resource.
	if err == nil {
		if typeCount < lifeCycle.MaxCount {
			logger.Debug("Adding new dynamic resources...")
			res := newResourceFromNewDynamicResourceLifeCycle(r.Storage.generateName(), r.now(), lifeCycle)
			if err := r.Storage.AddResource(res, tenant); err != nil {
				logger.WithError(err).Warningf("unable to add a new resource of type %s", rType)
			}
			logger.Infof("Added dynamic resource %s of type %s", res.Name, res.Type)
		}
	} else {
		logrus.WithError(err).Debug("Failed listing DRLC")
	}
}

// AcquireByState checks out resources of a given type without an owner,
// that matches a list of resources names.
// In: state - current state of the requested resource
//     dest - destination state of the requested resource
//     owner - requester of the resource
//     names - names of resource to acquire
// Out: A valid list of Resource object on success, or
//      ResourceNotFound error if target type resource does not exist in target state.
func (r *Ranch) AcquireByState(state, dest, owner string, names []string, tenant common.Tenant) ([]*common.Resource, error) {
	if names == nil {
		return nil, fmt.Errorf("must provide names of expected resources")
	}

	var returnRes []*common.Resource
	if err := retryOnConflict(func() error {
		rNames := util.NewSet(names...)

		allResources, err := r.Storage.GetResources(tenant)
		if err != nil {
			logrus.WithError(err).Errorf("could not get resources")
			return &ResourceNotFound{state}
		}

		var resources []*common.Resource

		for _, res := range allResources {

			if state != res.State || res.Owner != "" || !rNames.Has(res.Name) {
				continue
			}

			res.Owner = owner
			res.State = dest
			updatedRes, err := r.Storage.UpdateResource(res, tenant)
			if err != nil {
				return err
			}
			resources = append(resources, updatedRes)
			rNames.Delete(res.Name)
		}

		if rNames.Len() != 0 {
			missingResources := rNames.List()
			err := &ResourceNotFound{state}
			logrus.WithError(err).Errorf("could not find required resources %s", strings.Join(missingResources, ", "))
			returnRes = resources
			return err
		}
		returnRes = resources
		return nil
	}); err != nil {
		logrus.WithError(err).Error("AcquireByState failed")
		// Not a bug, we return what we got even on error.
		return returnRes, err
	}

	return returnRes, nil
}

// Release unsets owner for target resource and move it to a new state.
// In: name - name of the target resource
//     dest - destination state of the resource
//     owner - owner of the resource
// Out: nil on success, or
//      OwnerNotMatch error if owner does not match current owner of the resource, or
//      ResourceNotFound error if target named resource does not exist.
func (r *Ranch) Release(name, dest, owner string, tenant common.Tenant) error {
	logrus.Infof("Release")
	if err := retryOnConflict(func() error {
		res, err := r.Storage.GetResource(name, tenant)
		if err != nil {
			logrus.WithError(err).Errorf("unable to release resource %s", name)
			return &ResourceNotFound{name}
		}
		if owner != res.Owner {
			return &OwnerNotMatch{request: owner, owner: res.Owner}
		}

		res.Owner = ""
		res.State = dest

		if lf, err := r.Storage.GetDynamicResourceLifeCycle(res.Type, tenant); err == nil {
			// Assuming error means not existing as the only way to differentiate would be to list
			// all resources and find the right one which is more costly.
			if lf.LifeSpan != nil {
				expirationTime := r.now().Add(*lf.LifeSpan)
				res.ExpirationDate = &expirationTime
			}
		} else {
			res.ExpirationDate = nil
		}

		if _, err := r.Storage.UpdateResource(res, tenant); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logrus.WithError(err).Error("Release failed")
		return err
	}

	return nil
}

// Update updates the timestamp of a target resource.
// In: name  - name of the target resource
//     state - current state of the resource
//     owner - current owner of the resource
// 	   info  - information on how to use the resource
// Out: nil on success, or
//      OwnerNotMatch error if owner does not match current owner of the resource, or
//      ResourceNotFound error if target named resource does not exist, or
//      StateNotMatch error if state does not match current state of the resource.
func (r *Ranch) Update(name, owner, state string, ud *common.UserData, tenant common.Tenant) error {
	if err := retryOnConflict(func() error {
		res, err := r.Storage.GetResource(name, tenant)
		if err != nil {
			logrus.WithError(err).Errorf("could not find resource %s for update", name)
			return &ResourceNotFound{name}
		}
		if owner != res.Owner {
			return &OwnerNotMatch{request: owner, owner: res.Owner}
		}
		if state != res.State {
			return &StateNotMatch{res.State, state}
		}
		if res.UserData == nil {
			res.UserData = &common.UserData{}
		}
		res.UserData.Update(ud)
		if _, err := r.Storage.UpdateResource(res, tenant); err != nil {
			return err
		}
		return nil
	}); err != nil {
		logrus.WithError(err).Error("Update failed")
		return err
	}

	return nil
}

// Reset unstucks a type of stale resource to a new state.
// In: rtype - type of the resource
//     state - current state of the resource
//     expire - duration before resource's last update
//     dest - destination state of expired resources
// Out: map of resource name - resource owner.
func (r *Ranch) Reset(rtype, state string, expire time.Duration, dest string, tenant common.Tenant) (map[string]string, error) {
	var ret map[string]string
	if err := retryOnConflict(func() error {
		ret = make(map[string]string)
		resources, err := r.Storage.GetResources(tenant)
		if err != nil {
			return err
		}

		for _, res := range resources {
			if rtype != res.Type || state != res.State || res.Owner == "" || r.now().Sub(res.LastUpdate) < expire {
				continue
			}

			ret[res.Name] = res.Owner
			res.Owner = ""
			res.State = dest
			if _, err := r.Storage.UpdateResource(res, tenant); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		logrus.WithError(err).Error("Reset failed")
		return nil, err
	}

	return ret, nil
}

// SyncConfig updates resource list from a file
func (r *Ranch) SyncConfig(cfg *[]byte, tenant common.Tenant) error {
	config, err := common.ParseConfig(cfg)
	if err != nil {
		return err
	}
	if err := common.ValidateConfig(config); err != nil {
		return err
	}
	return r.Storage.SyncResources(config, tenant)
}

// StartDynamicResourceUpdater starts a goroutine which periodically
// updates all dynamic resources.
func (r *Ranch) StartDynamicResourceUpdater(updatePeriod time.Duration) {
	if updatePeriod == 0 {
		return
	}
	go func() {
		updateTick := time.NewTicker(updatePeriod).C
		for {
			select {
			case <-updateTick:
				if err := r.Storage.UpdateAllDynamicResources(); err != nil {
					logrus.WithError(err).Error("UpdateAllDynamicResources failed")
				}
			}
		}
	}()
}

// StartRequestGC starts the GC of expired requests
func (r *Ranch) StartRequestGC(gcPeriod time.Duration) {
	r.requestMgr.StartGC(gcPeriod)
}

// Metric returns a metric object with metrics filled in
func (r *Ranch) Metric(rtype string, tenant common.Tenant) (common.Metric, error) {
	metric := common.NewMetric(rtype)

	resources, err := r.Storage.GetResources(tenant)
	if err != nil {
		logrus.WithError(err).Error("cannot find resources")
		return metric, &ResourceNotFound{rtype}
	}

	for _, res := range resources {
		if res.Type != rtype {
			continue
		}

		metric.Current[res.State]++
		metric.Owners[res.Owner]++
	}

	if len(metric.Current) == 0 && len(metric.Owners) == 0 {
		return metric, &ResourceNotFound{rtype}
	}

	return metric, nil
}

// AllMetrics returns a list of Metric objects for all resource types.
func (r *Ranch) AllMetrics(tenant common.Tenant) ([]common.Metric, error) {
	resources, err := r.Storage.GetResources(tenant)
	if err != nil {
		logrus.WithError(err).Error("cannot get resources")
		return nil, err
	}

	metrics := map[string]common.Metric{}

	for _, res := range resources {
		metric, ok := metrics[res.Type]
		if !ok {
			metric = common.NewMetric(res.Type)
			metrics[res.Type] = metric
		}

		metric.Current[res.State]++
		metric.Owners[res.Owner]++
	}

	result := make([]common.Metric, 0, len(metrics))
	for _, metric := range metrics {
		result = append(result, metric)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Type < result[j].Type
	})
	return result, nil
}

// newResourceFromNewDynamicResourceLifeCycle creates a resource from DynamicResourceLifeCycle given a name and a time.
// Using this method helps make sure all the resources are created the same way.
//func newResourceFromNewDynamicResourceLifeCycle(name string, dlrc *crds.DRLCObject, now time.Time) *common.Resource {
//	//return crds.NewResource(name, dlrc.Name, dlrc.Spec.InitialState, "", now)
//	return nil
//}

func newResourceFromNewDynamicResourceLifeCycle(name string, now time.Time, lifeCycle common.DynamicResourceLifeCycle) *common.Resource {
	expirationTime := now.Add(*lifeCycle.LifeSpan)
	res := common.Resource{
		Type:           lifeCycle.Type,
		Name:           name,
		State:          lifeCycle.InitialState,
		LastUpdate:     now,
		ExpirationDate: &expirationTime,
	}
	return &res
}

func retryOnConflict(fn func() error) error {
	return fn()
}

func (r *Ranch) ValidateAuthToken(token, project string) (common.Tenant, error) {
	return r.Storage.ValidateToken(token, project)
}
