package persistence

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/cheld/miniprow/pkg/boskos/common"
	"google.golang.org/api/option"
)

type store struct {
	lock   sync.RWMutex
	client *firestore.Client
	ctx    context.Context
}

// NewResourceMemoryStorage creates an in memory persistence layer
func NewFirestore() (ResourcePersistence, TenantPersistence) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("../firestore.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	firestore := &store{
		client: client,
		ctx:    ctx,
	}
	return firestore, firestore
}

func (im *store) Close() {
	im.client.Close()
}

func (im *store) Add(r common.Resource, tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	base := im.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	_, err := base.Collection("resources").Doc(r.Name).Set(im.ctx, r)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	return nil
}

func (im *store) Delete(name string, tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()

	return nil
}

func (im *store) Update(r common.Resource, tenant common.Tenant) (common.Resource, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	return r, nil
}

func (im *store) Get(name string, tenant common.Tenant) (common.Resource, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	base := im.client.Collection("organizations").Doc(tenant.Organization).Collection("projects").Doc(tenant.Project)
	dsnap, err := base.Collection("resources").Doc(name).Get(im.ctx)
	if err != nil {
		return common.Resource{}, err
	}
	var r common.Resource
	dsnap.DataTo(&r)
	return r, nil
}

func (im *store) List(tenant common.Tenant) ([]common.Resource, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()

	return []common.Resource{}, nil
}

func (im *store) AddToken(token string, tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	_, err := im.client.Collection("tokens").Doc(token).Set(im.ctx, map[string]interface{}{
		"organization": tenant.Organization,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
	return nil
}

func (im *store) DeleteToken(tenant common.Tenant) error {
	im.lock.Lock()
	defer im.lock.Unlock()
	return nil
}

func (im *store) GetTenant(token, project string) (common.Tenant, error) {
	im.lock.Lock()
	defer im.lock.Unlock()
	return common.NewTenant(), nil
}
