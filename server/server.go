package server

import (
	"fmt"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/config"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
)

type Server struct {
	doppler  string
	port     int
	user     string
	pass     string
	listener net.Listener
}

func NewServer(c *config.ServerConfig) *Server {
	return &Server{
		doppler: c.Doppler,
		port:    c.Port,
		user:    c.User,
		pass:    c.Pass,
	}
}

func (s *Server) Start() {

	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	Logger.Info("Starting server at " + addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		Logger.Error("Failed-to-start-listener", err)
		os.Exit(1)
	}
	s.listener = listener

	s.registerHandlers()
	http.Serve(listener, nil)
}

func (s *Server) Stop() {
	s.listener.Close()
}

func (s *Server) registerHandlers() {
	handler := NewHandler(s.doppler)
	r := mux.NewRouter()
	r.Methods("GET").Path("/v1/apps/{appid}/metrics/memory").HandlerFunc(handler.GetMemoryMetric)
	http.Handle("/", r)
}
