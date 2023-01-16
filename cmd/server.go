package cmd

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	userAPI "github.com/semka95/balance-service/api"
	invoiceStore "github.com/semka95/balance-service/invoice/repository"
	transferStore "github.com/semka95/balance-service/transfer/repository"
	userStore "github.com/semka95/balance-service/user/repository"
	userUcase "github.com/semka95/balance-service/user/usecase"
)

// RestServer represents rest server
type RestServer struct {
	logger *zap.Logger
	config *Config
}

// NewServer creates rest server
func NewServer(logger *zap.Logger, config *Config) RestServer {
	return RestServer{
		logger: logger,
		config: config,
	}
}

// RunServer runs rest server
func (s *RestServer) RunServer() {
	// init database
	db, err := sql.Open(s.config.DBDriver, s.config.DBSource)
	if err != nil {
		s.logger.Error("can't open database connection", zap.Error(err), zap.String("db driver", s.config.DBDriver), zap.String("db source", s.config.DBSource))
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		s.logger.Error("can't ping database", zap.Error(err), zap.String("db driver", s.config.DBDriver), zap.String("db source", s.config.DBSource))
		return
	}

	// init router
	uStore := userStore.New(db)
	uUcase := userUcase.New(uStore, db)
	tStore := transferStore.New(db)
	invStore := invoiceStore.New(db)
	api := userAPI.API{}
	router := api.NewRouter(uStore, uUcase, tStore, invStore, db)

	// init http server
	srv := &http.Server{
		Addr:        s.config.HTTPServerAddress,
		Handler:     router,
		ReadTimeout: time.Duration(s.config.ReadTimeout) * time.Second,
		IdleTimeout: time.Duration(s.config.IdleTimeout) * time.Second,
	}
	// run server
	s.logger.Info("starting http server", zap.String("address", s.config.HTTPServerAddress))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("can't start server", zap.Error(err), zap.String("server address", s.config.HTTPServerAddress))
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	s.logger.Info("received interrupt signal, stopping server")
	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeout); err != nil {
		s.logger.Error("can't shutdown http server", zap.Error(err))
	}
}
