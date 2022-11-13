package server

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Server struct {
	Server *http.Server
}

func NewServer(config Config, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:              config.Host + ":" + config.Port, // адрес
			Handler:           handler,                         // обработчик запросов
			MaxHeaderBytes:    1 << 20,                         // максимальный объем хедера(1 МБ),
			ReadHeaderTimeout: 10 * time.Second,                // ограничение по времени для чтения
			WriteTimeout:      10 * time.Second,                // ограничение по времени для записи
		},
	}
}

func (s *Server) Run() error {
	return s.Server.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
