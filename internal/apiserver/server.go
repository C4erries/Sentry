package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host         string        // например, "0.0.0.0"
	Port         int           // например, 8080
	ReadTimeout  time.Duration // например, 5 * time.Second
	WriteTimeout time.Duration // например, 10 * time.Second
	IdleTimeout  time.Duration // например, 120 * time.Second
	Mode         string        // gin.ReleaseMode / gin.DebugMode
}

type Server struct {
	cfg    *Config
	engine *gin.Engine
}

func New(cfg Config) *Server {
	gin.SetMode(cfg.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	srv := &Server{
		cfg:    &cfg,
		engine: r,
	}

	srv.registerRoutes()

	return srv
}

func (s *Server) registerRoutes() {
	api := s.engine.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			// public
			users.GET("", s.listUsers)
			//users.GET("/:id", s.getUser)

			// protected
			// users.Use(s.authMiddleware.Middleware())
			//users.POST("", s.createUser)
		}
	}

	// здоровье
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// Run запускает HTTP-сервер и обеспечивает graceful shutdown
func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// здесь можно логировать ошибку через s.logger
			fmt.Printf("Listen error: %v\n", err)
		}
	}()
	fmt.Printf("Server listening on %s\n", addr)

	// ловим сигналы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpSrv.Shutdown(ctx)
}

func (s *Server) listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "listUsers"})
}
