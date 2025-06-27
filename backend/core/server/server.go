package server

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/runnable"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
)

// TODO CX-??: Remove server

// Server -
type Server struct {
	log    logger.Logger
	store  storer.Storer
	runner runnable.Runnable
}

// NewServer -
func NewServer(l logger.Logger, s storer.Storer, r runnable.Runnable) (*Server, error) {

	svr := Server{
		log:    l,
		store:  s,
		runner: r,
	}

	err := svr.Init()
	if err != nil {
		return nil, err
	}

	return &svr, nil
}

// Init -
func (svr *Server) Init() error {

	// TODO: alerting, retries
	return svr.runner.Init(svr.store)
}

// Run -
func (svr *Server) Run(args map[string]any) error {
	svr.log.Warn("(core) running server")
	return svr.runner.Run(args)
}
