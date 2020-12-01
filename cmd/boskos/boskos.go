package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cheld/cicd-bot/pkg/handlers"
	"github.com/cheld/cicd-bot/pkg/ranch"
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

	storage := ranch.NewStorage()

	r, _ := ranch.NewRanch("", storage, defaultRequestTTL)

	boskos := &http.Server{
		Handler: handlers.NewBoskosHandler(r),
		Addr:    ":8080",
	}
	fmt.Println("test")
	boskos.ListenAndServe()
}
