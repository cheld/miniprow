package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cheld/cicd-bot/pkg/boskos/handlers"
	"github.com/cheld/cicd-bot/pkg/boskos/ranch"
	"github.com/cheld/cicd-bot/pkg/boskos/storage"
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

	r, err := ranch.NewRanch("boskos.yaml", storage, defaultRequestTTL)
	if err != nil {
		fmt.Println(err)
	}

	boskos := &http.Server{
		Handler: handlers.NewBoskosHandler(r),
		Addr:    ":8080",
	}

	err = boskos.ListenAndServe()
	fmt.Println(err)
}
