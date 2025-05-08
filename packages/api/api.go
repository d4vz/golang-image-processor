package api

import (
	"context"
	"net/http"
)

type ApiConfig struct {
	port string
}

type Api struct {
	server *http.Server
	cfg    *ApiConfig
	mux    *http.ServeMux
}

func NewApiConfig(port string) *ApiConfig {
	return &ApiConfig{
		port: port,
	}
}

func NewApi(cfg *ApiConfig) *Api {
	mux := http.NewServeMux()
	return &Api{
		cfg: cfg,
		mux: mux,
	}
}

func (a *Api) Start() error {
	a.server = &http.Server{
		Addr:    a.cfg.port,
		Handler: a.mux,
	}
	return a.server.ListenAndServe()
}

func (a *Api) Stop() error {
	if a.server == nil {
		return nil
	}
	return a.server.Shutdown(context.Background())
}

func (a *Api) Handle(path string, handler http.HandlerFunc) {
	a.mux.HandleFunc(path, handler)
}
