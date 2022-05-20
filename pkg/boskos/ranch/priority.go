/*
Copyright 2019 The Kubernetes Authors.

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
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// request stores request information with expiration
type request struct {
	id         string
	expiration time.Time
	// Used to calculate since when this resource has been acquired
	createdAt time.Time
}

type requestNode struct {
	id   string
	next *requestNode
}

type requestLinkedList struct {
	start, end *requestNode
}

func (l *requestLinkedList) Append(id string) {
	if l.start == nil {
		l.start = &requestNode{id: id}
		l.end = l.start
		return
	}
	l.end.next = &requestNode{id: id}
	l.end = l.end.next
}

func (l *requestLinkedList) Delete(requestID string) {
	previous := l.start
	for n := l.start; n != nil; n = n.next {
		if n.id == requestID {
			if n == l.start {
				l.start = l.start.next
			}
			if n == l.end {
				l.end = previous
			} else {
				previous.next = n.next
			}
			return
		}
		previous = n
	}
}

func (l *requestLinkedList) Range(f func(string) bool) {
	for n := l.start; n != nil; n = n.next {
		if b := f(n.id); !b {
			break
		}
	}
}

// requestQueue is a simple FIFO queue for requests.
type requestQueue struct {
	lock        sync.RWMutex
	requestList *requestLinkedList
	requestMap  map[string]request
}

func newRequestQueue() *requestQueue {
	return &requestQueue{
		requestMap:  map[string]request{},
		requestList: &requestLinkedList{},
	}
}

// update updates expiration time is updated if already present,
// add a new requestID at the end otherwise (FIFO)
func (rq *requestQueue) update(requestID string, newExpiration, now time.Time) bool {
	rq.lock.Lock()
	defer rq.lock.Unlock()
	req, exists := rq.requestMap[requestID]
	if !exists {
		req = request{id: requestID, createdAt: now}
		rq.requestList.Append(requestID)
		logrus.Infof("request id %s added", requestID)
	}
	// Update timestamp
	req.expiration = newExpiration
	rq.requestMap[requestID] = req
	logrus.Infof("request id %s set to expire at %v", requestID, newExpiration)
	return !exists
}

// delete an element
func (rq *requestQueue) delete(requestID string) {
	rq.lock.Lock()
	defer rq.lock.Unlock()
	delete(rq.requestMap, requestID)
	rq.requestList.Delete(requestID)
}

// cleanup checks for all expired  or marked for deletion items and delete them.
func (rq *requestQueue) cleanup(now time.Time) {
	rq.lock.Lock()
	defer rq.lock.Unlock()
	newRequestList := &requestLinkedList{}
	newRequestMap := map[string]request{}
	rq.requestList.Range(func(requestID string) bool {
		req := rq.requestMap[requestID]
		// Checking expiration
		if now.After(req.expiration) {
			logrus.Infof("request id %s expired", req.id)
			return true
		}
		// Keeping
		newRequestList.Append(requestID)
		newRequestMap[requestID] = req
		return true
	})
	rq.requestMap = newRequestMap
	rq.requestList = newRequestList
}

// getRank provides the rank of a given requestID following the order it was added (FIFO).
// If requestID is an empty string, getRank assumes it is added last (lowest rank + 1).
func (rq *requestQueue) getRank(requestID string, ttl time.Duration, now time.Time) (int, bool) {
	// not considering empty requestID as new
	var new bool
	if requestID != "" {
		new = rq.update(requestID, now.Add(ttl), now)
	}
	rank := 1
	rq.lock.RLock()
	defer rq.lock.RUnlock()
	rq.requestList.Range(func(existingID string) bool {
		req := rq.requestMap[existingID]
		if now.After(req.expiration) {
			logrus.Infof("request id %s expired", req.id)
			return true
		}
		if requestID == existingID {
			return false
		}
		rank++
		return true
	})
	return rank, new
}

func (rq *requestQueue) isEmpty() bool {
	rq.lock.Lock()
	defer rq.lock.Unlock()
	return len(rq.requestMap) == 0
}

// RequestManager facilitates management of RequestQueues for a set of (resource type, resource state) tuple.
type RequestManager struct {
	lock     sync.Mutex
	requests map[interface{}]*requestQueue
	ttl      time.Duration
	stopGC   context.CancelFunc
	wg       sync.WaitGroup
	// For testing only
	now func() time.Time
}

// NewRequestManager creates a new RequestManager
func NewRequestManager(ttl time.Duration) *RequestManager {
	return &RequestManager{
		requests: map[interface{}]*requestQueue{},
		ttl:      ttl,
		now:      time.Now,
	}
}

func (rm *RequestManager) cleanup(now time.Time) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	for key, rq := range rm.requests {
		logrus.Infof("cleaning up %v request queue", key)
		rq.cleanup(now)
		if rq.isEmpty() {
			delete(rm.requests, key)
		}
	}
}

// StartGC starts a goroutine that will call cleanup every gcInterval
func (rm *RequestManager) StartGC(gcPeriod time.Duration) {
	ctx, stop := context.WithCancel(context.Background())
	rm.stopGC = stop
	tick := time.Tick(gcPeriod)
	rm.wg.Add(1)
	go func() {
		logrus.Info("starting cleanup go routine")
		defer logrus.Info("exiting cleanup go routine")
		defer rm.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				rm.cleanup(rm.now())
			}
		}

	}()
}

// StopGC is a blocking call that will stop the GC goroutine.
func (rm *RequestManager) StopGC() {
	if rm.stopGC != nil {
		rm.stopGC()
		rm.wg.Wait()
	}
}

// GetRank provides the rank of a given request and whether request is new (was added)
func (rm *RequestManager) GetRank(key interface{}, id string) (int, bool) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rq := rm.requests[key]
	if rq == nil {
		rq = newRequestQueue()
		rm.requests[key] = rq
	}
	return rq.getRank(id, rm.ttl, rm.now())
}

// GetCreatedAt returns when the request was created
func (rm *RequestManager) GetCreatedAt(key interface{}, id string) (time.Time, error) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	var createdTime time.Time
	rq := rm.requests[key]
	if rq == nil {
		//This should never happen
		return createdTime, fmt.Errorf("failed to get the created time because the request queue is nil for the key %v", key)
	}

	req, exists := rq.requestMap[id]
	if !exists {
		//This should never happen
		return createdTime, fmt.Errorf("failed to get the created time because the request does not exist for the id %s", id)
	}
	return req.createdAt, nil
}

// Delete deletes a specific request such that it is not accounted in the next GetRank call.
// This is usually called when the request has been fulfilled.
func (rm *RequestManager) Delete(key interface{}, requestID string) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rq := rm.requests[key]
	if rq != nil {
		rq.delete(requestID)
	}
}
