package server

import (
	"context"
	bannerrepo "github.com/SanExpett/banners-backend/internal/banner/repository"
	bannerusecases "github.com/SanExpett/banners-backend/internal/banner/usecases"
	"github.com/SanExpett/banners-backend/internal/server/delivery/mux"
	"github.com/SanExpett/banners-backend/internal/server/repository"
	userrepo "github.com/SanExpett/banners-backend/internal/user/repository"
	userusecases "github.com/SanExpett/banners-backend/internal/user/usecases"
	"github.com/SanExpett/banners-backend/pkg/config"
	"github.com/SanExpett/banners-backend/pkg/my_logger"
	"net/http"
	"strings"
	"time"
)

const (
	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(config *config.Config) error {
	baseCtx := context.Background()

	pool, err := repository.NewPgxPool(baseCtx, config.URLDataBase)
	if err != nil {
		return err //nolint:wrapcheck
	}

	logger, err := my_logger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer logger.Sync()

	userStorage, err := userrepo.NewUserStorage(pool)
	if err != nil {
		return err
	}

	userService, err := userusecases.NewUserService(userStorage)
	if err != nil {
		return err
	}

	bannerStorage, err := bannerrepo.NewBannerStorage(pool)
	if err != nil {
		return err
	}

	bannerService, err := bannerusecases.NewBannerService(bannerStorage)
	if err != nil {
		return err
	}

	handler, err := mux.NewMux(baseCtx, mux.NewConfigMux(config.AllowOrigin,
		config.Schema, config.PortServer), userService, bannerService, logger)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.PortServer,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    basicTimeout,
		WriteTimeout:   basicTimeout,
	}

	logger.Infof("Start server:%s", config.PortServer)

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}
