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
	store := persistence.NewFirestore()
	r := common.NewResource("someresource", "free", "used", "cheld", time.Now())
	store.Add(r, common.NewTenant())
	store.AddToken("aaaaaaaaa", common.NewTenant())
	l, _ := store.Get("someresource", common.NewTenant())
	fmt.Println(l)
}
