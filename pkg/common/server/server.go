package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cheld/miniprow/pkg/common/core"
	"github.com/cheld/miniprow/pkg/common/info"
	"github.com/cheld/miniprow/pkg/common/notification"
)

type CommonServer struct {
	mux *http.ServeMux
}

func NewHandler(notifyer *notification.Dispatcher) *CommonServer {
	server := CommonServer{
		mux: http.NewServeMux(),
	}
	server.mux.Handle("/health", handleHealth())
	server.mux.Handle("/version", handleVersion())
	notifyer.Register(func(*core.Event, core.Tenant, context.Context) {
		fmt.Println("----------------event received")
	}, "github_comment")
	return &server
}

func (s *CommonServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func handleHealth() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Ok\n")
	}
}

func handleVersion() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Version %v\n", info.Version)
		fmt.Fprintf(res, "Commit %v\n", info.Commit)
	}
}
