package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cheld/cicd-bot/pkg/boskos/handlers"
	"github.com/cheld/cicd-bot/pkg/boskos/ranch"
	"github.com/cheld/cicd-bot/pkg/boskos/storage"
	"github.com/sirupsen/logrus"
)

const (
	defaultDynamicResourceUpdatePeriod = 10 * time.Minute
	defaultRequestTTL                  = 30 * time.Second
	defaultRequestGCPeriod             = time.Minute
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	// collect data on mutex holders and blocking profiles
	//runtime.SetBlockProfileRate(1)
	//runtime.SetMutexProfileFraction(1)

	storage := ranch.NewStorage(storage.NewMemoryStorage())

	r, _ := ranch.NewRanch("boskos.yaml", storage, defaultRequestTTL)

	boskos := &http.Server{
		Handler: handlers.NewBoskosHandler(r),
		Addr:    ":8080",
	}
	fmt.Println("test")
	boskos.ListenAndServe()
}
