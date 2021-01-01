package server

import (
	"fmt"
	"net/http"

	"github.com/cheld/miniprow/pkg/common/info"
)

type CommonServer struct {
	mux *http.ServeMux
}

func NewHandler() *CommonServer {
	server := CommonServer{
		mux: http.NewServeMux(),
	}
	server.mux.Handle("/common/health", handleHealth())
	server.mux.Handle("/common/version", handleVersion())
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
	}
}
