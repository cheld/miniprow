package persistence

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/cheld/miniprow/pkg/boskos/common"
	"github.com/cheld/miniprow/pkg/common/util"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type store struct {
	lock   sync.RWMutex
	client *firestore.Client
	ctx    context.Context
}

// NewResourceMemoryStorage creates an in memory persistence layer
func NewFirestore() Persistence {
	// Use a service account
	ctx := context.Background()
	var app *firebase.App
	var err error
	if util.FileExists("../firestore.json") {
		sa := option.WithCredentialsFile("../firestore.json")
		app, err = firebase.NewApp(ctx, nil, sa)
	} else {
		conf := &firebase.Config{ProjectID: "smart-altar-272110"}
		app, err = firebase.NewApp(ctx, conf)
	}
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	store := &store{
		client: client,
		ctx:    ctx,
	}
	return store
}

func (s *store) Close() {
	s.client.Close()
}

func (s *store) Add(r common.Resource, tenant common.Tenant) error {
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	_, err := baseQuery.Collection("resources").Doc(r.Name).Set(s.ctx, r)
	if err != nil {
		log.Fatalf("Failed adding resource to firestore: %v", err)
	}
	return nil
}

func (s *store) Delete(name string, tenant common.Tenant) error {
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	_, err := baseQuery.Collection("resources").Doc(name).Delete(s.ctx)
	if err != nil {
		log.Fatalf("Failed deleting resource from firestore: %v", err)
	}
	return nil
}

func (s *store) Update(r common.Resource, tenant common.Tenant) (common.Resource, error) {
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	_, err := baseQuery.Collection("resources").Doc(r.Name).Set(s.ctx, r)
	if err != nil {
		log.Fatalf("Failed updating resource to firestore: %v", err)
	}
	return r, nil
}

func (s *store) Get(name string, tenant common.Tenant) (common.Resource, error) {
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	dsnap, err := baseQuery.Collection("resources").Doc(name).Get(s.ctx)
	if err != nil {
		return common.Resource{}, err
	}
	var r common.Resource
	dsnap.DataTo(&r)
	return r, nil
}

func (s *store) List(tenant common.Tenant) ([]common.Resource, error) {
	result := []common.Resource{}
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	iter := baseQuery.Collection("resources").Documents(s.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return result, err
		}
		var r common.Resource
		doc.DataTo(&r)
		result = append(result, r)
	}
	return result, nil
}

func (s *store) AddDynamicResourceLifeCycle(r common.DynamicResourceLifeCycle, tenant common.Tenant) error {
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	_, err := baseQuery.Collection("drlc").Doc(r.Type).Set(s.ctx, r)
	if err != nil {
		log.Fatalf("Failed adding resource to firestore: %v", err)
	}
	return nil
}

func (s *store) GetDynamicResourceLifeCycle(rtype string, tenant common.Tenant) (common.DynamicResourceLifeCycle, error) {
	baseQuery := s.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	dsnap, err := baseQuery.Collection("resources").Doc(rtype).Get(s.ctx)
	if err != nil {
		return common.DynamicResourceLifeCycle{}, err
	}
	var r common.DynamicResourceLifeCycle
	dsnap.DataTo(&r)
	return r, nil
}

func (s *store) AddToken(token string, tenant common.Tenant) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, err := s.client.Collection("tokens").Doc(token).Set(s.ctx, map[string]interface{}{
		"organization": tenant.Organization,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	return nil
}

func (s *store) DeleteToken(tenant common.Tenant) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return nil
}

func (s *store) GetTenantFromToken(token, project string) (common.Tenant, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return common.NewTenant(), nil
}
