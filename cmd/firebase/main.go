package main

import (
	"fmt"
	"time"

	"github.com/cheld/miniprow/pkg/boskos/common"
	"github.com/cheld/miniprow/pkg/boskos/persistence"
)

func main() {
	fmt.Println("hello")
	// Use a service account
	resources, tenants := persistence.NewFirestore()
	r := common.NewResource("someresource", "free", "used", "cheld", time.Now())
	resources.Add(r, common.NewTenant())
	tenants.AddToken("aaaaaaaaa", common.NewTenant())
	l, _ := resources.Get("someresource", common.NewTenant())
	fmt.Println(l)
}
