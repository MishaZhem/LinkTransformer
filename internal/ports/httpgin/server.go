package httpgin

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
}

func NewHTTPServer(port string) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{}
	s.server = &http.Server{
		Addr:    port,
		Handler: s.Handler(),
	}
	return s
}

func CustomLogger(c *gin.Context) {
	t := time.Now()

	c.Next()

	latency := time.Since(t)
	status := c.Writer.Status()

	log.Println("latency", latency, "method", c.Request.Method, "path", c.Request.URL.Path, "status", status)
}

func (s *Server) Listen() error {
	return s.server.ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	a := gin.New()
	a.Use(CustomLogger)
	a.Use(gin.Recovery())
	AppRouter(&a.RouterGroup)
	return a
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
