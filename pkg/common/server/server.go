package server

import (
	"fmt"
	"net/http"

	"github.com/cheld/miniprow/pkg/common/config"
)

func Register(mux *http.ServeMux) {
	mux.Handle("/health", handleHealth())
	mux.Handle("/version", handleVersion())
}

func handleHealth() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Ok\n")
	}
}

func handleVersion() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Version %v\n", config.Version)
	}
}
