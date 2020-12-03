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
	"time"

	"github.com/cheld/cicd-bot/pkg/boskos/common"

	"github.com/cheld/cicd-bot/pkg/boskos/storage"
)

// Storage is used to decouple ranch functionality with the resource persistence layer
type Storage struct {
	persistence storage.PersistenceLayer
	// For testing
	now          func() time.Time
	generateName func() string
}

// NewStorage instantiates a new Storage with a PersistenceLayer implementation
func NewStorage(persistence storage.PersistenceLayer) *Storage {
	return &Storage{
		persistence:  persistence,
		now:          func() time.Time { return time.Now() },
		generateName: common.GenerateDynamicResourceName,
	}
}

// AddResource adds a new resource
func (s *Storage) AddResource(resource *common.Resource) error {
	return s.persistence.Add(*resource)
}

// DeleteResource deletes a resource if it exists, errors otherwise
func (s *Storage) DeleteResource(name string) error {
	return s.persistence.Delete(name)
}

// UpdateResource updates a resource if it exists, errors otherwise
func (s *Storage) UpdateResource(resource *common.Resource) (*common.Resource, error) {
	resource.LastUpdate = s.now()
	s.persistence.Update(*resource)
	return resource, nil
}

// GetResource gets an existing resource, errors otherwise
func (s *Storage) GetResource(name string) (*common.Resource, error) {
	res, err := s.persistence.Get(name)
	return &res, err
}

// GetResources list all resources
func (s *Storage) GetResources() ([]*common.Resource, error) {
	r, _ := s.persistence.List()
	resourceList := make([]*common.Resource, len(r))
	for i := 0; i < len(r); i++ {
		resourceList[i] = &r[i]
	}
	return resourceList, nil
}

// AddDynamicResourceLifeCycle adds a new dynamic resource life cycle
func (s *Storage) AddDynamicResourceLifeCycle(resource string) error {
	// resource.Namespace = s.namespace
	// return s.client.Create(s.ctx, resource)
	return nil
}

// DeleteDynamicResourceLifeCycle deletes a dynamic resource life cycle if it exists, errors otherwise
func (s *Storage) DeleteDynamicResourceLifeCycle(name string) error {
	// o := &crds.DRLCObject{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      name,
	// 		Namespace: s.namespace,
	// 	},
	// }
	// return s.client.Delete(s.ctx, o)
	return nil
}

// UpdateDynamicResourceLifeCycle updates a dynamic resource life cycle. if it exists, errors otherwise
func (s *Storage) UpdateDynamicResourceLifeCycle(resource string) (string, error) {
	// resource.Namespace = s.namespace
	// if err := s.client.Update(s.ctx, resource); err != nil {
	// 	return nil, fmt.Errorf("failed to update dlrc %s: %w", resource.Name, err)
	// }

	return "", nil
}

// GetDynamicResourceLifeCycle gets an existing dynamic resource life cycle, errors otherwise
func (s *Storage) GetDynamicResourceLifeCycle(name string) (string, error) {
	// drlc := &crds.DRLCObject{}
	// nn := types.NamespacedName{Namespace: s.namespace, Name: name}
	// if err := s.client.Get(s.ctx, nn, drlc); err != nil {
	// 	return nil, fmt.Errorf("failed to get dlrc %s: %q", name, err)
	// }

	return "", nil
}

// GetDynamicResourceLifeCycles list all dynamic resource life cycle
func (s *Storage) GetDynamicResourceLifeCycles() (string, error) {
	// drlcList := &crds.DRLCObjectList{}
	// if err := s.client.List(s.ctx, drlcList, ctrlruntimeclient.InNamespace(s.namespace)); err != nil {
	// 	return nil, fmt.Errorf("failed to list drlcs: %v", err)
	// }

	return "", nil
}

// SyncResources will update static and dynamic resources.
// It will add new resources to storage and try to remove newly deleted resources
// from storage.
// If the newly deleted resource is currently held by a user, the deletion will
// yield to next update cycle.
func (s *Storage) SyncResources(config *common.BoskosConfig) error {
	if config == nil {
		return nil
	}

	for _, entry := range config.Resources {
		if entry.IsDRLC() {
			// 			newDRLCByType[entry.Type] = *crds.FromDynamicResourceLifecycle(common.NewDynamicResourceLifeCycleFromConfig(entry))
		} else {
			for _, res := range common.NewResourcesFromConfig(entry) {
				s.persistence.Add(res)
			}
		}
	}

	// 	if err := func() error {
	// 		s.resourcesLock.Lock()
	// 		defer s.resourcesLock.Unlock()

	// 		resources, err := s.GetResources()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		existingDRLC, err := s.GetDynamicResourceLifeCycles()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		for _, dRLC := range existingDRLC.Items {
	// 			existingDRLCByType[dRLC.Name] = dRLC
	// 		}

	// 		// Split resources between static and dynamic resources
	// 		lifeCycleTypes := sets.String{}

	// 		for _, lc := range existingDRLC.Items {
	// 			lifeCycleTypes.Insert(lc.Name)
	// 		}
	// 		// Considering the migration case from mason resources to dynamic resources.
	// 		// Dynamic resources already exist but they don't have an associated DRLC
	// 		for _, lc := range newDRLCByType {
	// 			lifeCycleTypes.Insert(lc.Name)
	// 		}

	// 		for _, res := range resources.Items {
	// 			if !lifeCycleTypes.Has(res.Spec.Type) {
	// 				existingSRByName[res.Name] = res
	// 			}
	// 		}

	// 		if err := s.syncStaticResources(newSRByName, existingSRByName); err != nil {
	// 			return err
	// 		}
	// 		if err := s.syncDynamicResourceLifeCycles(newDRLCByType, existingDRLCByType); err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}(); err != nil {
	// 		return err
	// 	}

	// if err := s.UpdateAllDynamicResources(); err != nil {
	// 	logrus.WithError(err).Error("Encountered error updating dynamic resources.")
	// 	return err
	// }

	return nil
}

// updateDynamicResources updates dynamic resource based on an existing dynamic resource life cycle.
// It will make sure than MinCount of resource exists, and attempt to delete expired and resources over MaxCount.
// If resources are held by another user than Boskos, they will be deleted in a following cycle.
// func (s *Storage) updateDynamicResources(lifecycle *crds.DRLCObject, resources []crds.ResourceObject) (toAdd, toDelete []crds.ResourceObject) {
// 	var notInUseRes []crds.ResourceObject
// 	tombStoned := 0
// 	toBeDeleted := 0
// 	for _, r := range resources {
// 		if r.Status.Owner != "" {
// 			// We can only delete resources not in use.
// 			continue
// 		}
// 		if r.Status.State == common.Tombstone {
// 			// Ready to be fully deleted.
// 			toDelete = append(toDelete, r)
// 			tombStoned++
// 		} else if r.Status.State == common.ToBeDeleted {
// 			// Already in the process of cleaning up.
// 			// Don't create new resources yet, but also don't delete additional
// 			// resources.
// 			toBeDeleted++
// 		} else {
// 			if r.Status.ExpirationDate != nil && s.now().After(*r.Status.ExpirationDate) {
// 				// Expired. Don't decrement the active count until it's tombstoned,
// 				// however, as it might be depending on other resources that need
// 				// to be released first.
// 				toDelete = append(toDelete, r)
// 				toBeDeleted++
// 			} else {
// 				notInUseRes = append(notInUseRes, r)
// 			}
// 		}
// 	}

// 	// Tombstoned resources are ready to be fully deleted, so replace them if necessary.
// 	activeCount := len(resources) - tombStoned
// 	for i := activeCount; i < lifecycle.Spec.MinCount; i++ {
// 		res := newResourceFromNewDynamicResourceLifeCycle(s.generateName(), lifecycle, s.now())
// 		toAdd = append(toAdd, *res)
// 		activeCount++
// 	}

// 	// ToBeDeleted resources may take some time to be fully cleaned up.
// 	// We can temporarily exceed MaxCount while these are being cleaned up,
// 	// particularly if MaxCount was recently lowered.
// 	numberOfResToDelete := activeCount - toBeDeleted - lifecycle.Spec.MaxCount
// 	// Sorting to get consistent deletion mechanism (ease testing)
// 	sort.SliceStable(notInUseRes, func(i, j int) bool {
// 		return notInUseRes[i].Name > notInUseRes[j].Name
// 	})
// 	for i := 0; i < len(notInUseRes); i++ {
// 		res := notInUseRes[i]
// 		if i < numberOfResToDelete {
// 			toDelete = append(toDelete, res)
// 		}
// 	}
// 	logrus.Infof("DRLC type %s: adding %+v, deleting %+v", lifecycle.Name, toAdd, toDelete)
// 	return
// }

// UpdateAllDynamicResources queries for all existing DynamicResourceLifeCycles
// and dynamic resources and calls updateDynamicResources for each type.
// This ensures that the MinCount and MaxCount parameters are honored, that
// any expired resources are deleted, and that any Tombstoned resources are
// completely removed.
func (s *Storage) UpdateAllDynamicResources() error {
	// s.resourcesLock.Lock()
	// defer s.resourcesLock.Unlock()

	// return retryOnConflict(retry.DefaultBackoff, func() error {
	// 	var resToAdd, resToDelete []crds.ResourceObject
	// 	var dRLCToDelete []crds.DRLCObject
	// 	existingDRLCByType := map[string]crds.DRLCObject{}
	// 	existingDRsByType := map[string][]crds.ResourceObject{}

	// 	resources, err := s.GetResources()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	existingDRLC, err := s.GetDynamicResourceLifeCycles()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, dRLC := range existingDRLC.Items {
	// 		existingDRLCByType[dRLC.Name] = dRLC
	// 	}

	// 	// Filter to only look at dynamic resources
	// 	for _, res := range resources.Items {
	// 		if _, ok := existingDRLCByType[res.Spec.Type]; ok {
	// 			existingDRsByType[res.Spec.Type] = append(existingDRsByType[res.Spec.Type], res)
	// 		}
	// 	}

	// 	for resType, dRLC := range existingDRLCByType {
	// 		existingDRs := existingDRsByType[resType]
	// 		toAdd, toDelete := s.updateDynamicResources(&dRLC, existingDRs)
	// 		resToAdd = append(resToAdd, toAdd...)
	// 		resToDelete = append(resToDelete, toDelete...)

	// 		if dRLC.Spec.MinCount == 0 && dRLC.Spec.MaxCount == 0 {
	// 			currentCount := len(existingDRs)
	// 			addCount := len(resToAdd)
	// 			delCount := len(resToDelete)
	// 			if addCount == 0 && (currentCount == 0 || currentCount == delCount) {
	// 				dRLCToDelete = append(dRLCToDelete, dRLC)
	// 			}
	// 		}
	// 	}

	// 	if err := s.persistResources(resToAdd, resToDelete, true); err != nil {
	// 		return err
	// 	}

	// 	if len(dRLCToDelete) > 0 {
	// 		if err := s.persistDynamicResourceLifeCycles(nil, nil, dRLCToDelete); err != nil {
	// 			return err
	// 		}
	// 	}

	// 	return nil
	// })
	return nil
}

// syncDynamicResourceLifeCycles compares the new DRLC configuration against
// the current configuration. If a DRLC has been deleted from the new
// configuration, it is updated to indicate that its dynamic resources should
// be removed.
// No dynamic resources are created, deleted, or modified by this function.
// func (s *Storage) syncDynamicResourceLifeCycles(newDRLCByType, existingDRLCByType map[string]crds.DRLCObject) error {
// 	var dRLCToUpdate, dRLCToAdd []crds.DRLCObject

// 	for _, existingDRLC := range existingDRLCByType {
// 		newDRLC, existsInNew := newDRLCByType[existingDRLC.Name]
// 		if existsInNew {
// 			// Copy the ObjectMeta so we can compare, the update doesn't fail due to unset ResourceVersion
// 			// and to keep any additional metadata that was added.
// 			newDRLC.ObjectMeta = *existingDRLC.ObjectMeta.DeepCopy()
// 			if !reflect.DeepEqual(existingDRLC, newDRLC) {
// 				dRLCToUpdate = append(dRLCToUpdate, newDRLC)
// 			}
// 		} else {
// 			// Mark for deletion of all associated dynamic resources.
// 			existingDRLC.Spec.MinCount = 0
// 			existingDRLC.Spec.MaxCount = 0
// 			dRLCToUpdate = append(dRLCToUpdate, existingDRLC)
// 		}
// 	}

// 	for _, newDRLC := range newDRLCByType {
// 		_, exists := existingDRLCByType[newDRLC.Name]
// 		if !exists {
// 			dRLCToAdd = append(dRLCToAdd, newDRLC)
// 		}
// 	}

// 	return s.persistDynamicResourceLifeCycles(dRLCToUpdate, dRLCToAdd, nil)
// }

// func (s *Storage) persistResources(resToAdd, resToDelete []crds.ResourceObject, dynamic bool) error {
// 	var errs []error
// 	for _, r := range resToDelete {
// 		// If currently busy, yield deletion to later cycles.
// 		if r.Status.Owner != "" {
// 			continue
// 		}
// 		if dynamic {
// 			// Only delete resource in tombsone state and mark the other as to deleted
// 			// This is necessary for dynamic resources that depends on other resources
// 			// as they need to be released to prevent leak.
// 			if r.Status.State == common.Tombstone {
// 				logrus.Infof("Deleting resource %s", r.Name)
// 				if err := s.DeleteResource(r.Name); err != nil {
// 					errs = append(errs, err)
// 				}
// 			} else if r.Status.State != common.ToBeDeleted {
// 				r.Status.State = common.ToBeDeleted
// 				logrus.Infof("Marking resource to be deleted %s", r.Name)
// 				if _, err := s.UpdateResource(&r); err != nil {
// 					errs = append(errs, err)
// 				}
// 			}
// 		} else {
// 			// Static resources can be deleted right away.
// 			logrus.Infof("Deleting resource %s", r.Name)
// 			if err := s.DeleteResource(r.Name); err != nil {
// 				errs = append(errs, err)
// 			}
// 		}
// 	}

// 	for _, r := range resToAdd {
// 		logrus.Infof("Adding resource %s", r.Name)
// 		r.Status.LastUpdate = s.now()
// 		if err := s.AddResource(&r); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}

// 	return utilerrors.NewAggregate(errs)
// }

// func (s *Storage) persistDynamicResourceLifeCycles(dRLCToUpdate, dRLCToAdd, dRLCToDelelete []crds.DRLCObject) error {
// 	var errs []error
// 	remainingTypes := map[string]bool{}
// 	updatedResources, err := s.GetResources()
// 	if err != nil {
// 		return err
// 	}
// 	for _, res := range updatedResources.Items {
// 		remainingTypes[res.Spec.Type] = true
// 	}

// 	for _, dRLC := range dRLCToDelelete {
// 		// Only delete a dynamic resource if all resources are gone
// 		if !remainingTypes[dRLC.Name] {
// 			logrus.Infof("Deleting resource type life cycle %s", dRLC.Name)
// 			if err := s.DeleteDynamicResourceLifeCycle(dRLC.Name); err != nil {
// 				errs = append(errs, err)
// 			}
// 		} else {
// 			// Mark this DRLC as pending deletion by setting min and max count to zero.
// 			dRLC.Spec.MinCount = 0
// 			dRLC.Spec.MaxCount = 0
// 			dRLCToUpdate = append(dRLCToUpdate, dRLC)
// 		}
// 	}

// 	for _, DRLC := range dRLCToAdd {
// 		logrus.Infof("Adding resource type life cycle %s", DRLC.Name)
// 		if err := s.AddDynamicResourceLifeCycle(&DRLC); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}

// 	for _, dRLC := range dRLCToUpdate {
// 		logrus.Infof("Updating resource type life cycle %s", dRLC.Name)
// 		if _, err := s.UpdateDynamicResourceLifeCycle(&dRLC); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}

// 	return utilerrors.NewAggregate(errs)
// }

// func (s *Storage) syncStaticResources(newResourcesByName, existingResourcesByName map[string]crds.ResourceObject) error {
// 	var resToAdd, resToDelete []crds.ResourceObject

// 	// Delete resources
// 	for _, res := range existingResourcesByName {
// 		_, exists := newResourcesByName[res.Name]
// 		if !exists {
// 			resToDelete = append(resToDelete, res)
// 		}
// 	}

// 	// Add new resources
// 	for _, res := range newResourcesByName {
// 		_, exists := existingResourcesByName[res.Name]
// 		if !exists {
// 			resToAdd = append(resToAdd, res)
// 		}
// 	}
// 	return s.persistResources(resToAdd, resToDelete, false)
// }
