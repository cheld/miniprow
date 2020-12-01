package main

import (
	"fmt"
	"net/http"

	"github.com/cheld/cicd-bot/pkg/handlers"
	"github.com/cheld/cicd-bot/pkg/ranch"
)

func main() {
	//logrusutil.ComponentInit()

	// collect data on mutex holders and blocking profiles
	//runtime.SetBlockProfileRate(1)
	//runtime.SetMutexProfileFraction(1)

	storage := ranch.NewStorage(nil, nil, nil)

	r, _ := ranch.NewRanch(nil, nil, nil)

	boskos := &http.Server{
		Handler: handlers.NewBoskosHandler(r),
		Addr:    ":8080",
	}
	fmt.Println("test")
}
