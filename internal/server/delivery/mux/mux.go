package mux

import (
	"context"
	"github.com/SanExpett/banners-backend/pkg/middleware"
	"net/http"

	bannerdelivery "github.com/SanExpett/banners-backend/internal/banner/delivery"
	userdelivery "github.com/SanExpett/banners-backend/internal/user/delivery"

	"go.uber.org/zap"
)

type ConfigMux struct {
	addrOrigin string
	schema     string
	portServer string
}

func NewConfigMux(addrOrigin string, schema string, portServer string) *ConfigMux {
	return &ConfigMux{
		addrOrigin: addrOrigin,
		schema:     schema,
		portServer: portServer,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux, userService userdelivery.IUserService,
	bannerService bannerdelivery.IBannerService, logger *zap.SugaredLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	userHandler, err := userdelivery.NewUserHandler(userService)
	if err != nil {
		return nil, err
	}

	bannerHandler, err := bannerdelivery.NewBannerHandler(bannerService)
	if err != nil {
		return nil, err
	}

	router.Handle("/api/v1/signup", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.SignUpHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/signin", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.SignInHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/logout", middleware.Context(ctx,
		middleware.SetupCORS(userHandler.LogOutHandler, configMux.addrOrigin, configMux.schema)))

	router.Handle("/api/v1/banner/add", middleware.Context(ctx,
		middleware.SetupCORS(bannerHandler.AddBannerHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/banner/get", middleware.Context(ctx,
		middleware.SetupCORS(bannerHandler.GetBannerHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/banner/delete", middleware.Context(ctx,
		middleware.SetupCORS(bannerHandler.DeleteBannerHandler, configMux.addrOrigin, configMux.schema)))
	router.Handle("/api/v1/banner/get_list", middleware.Context(ctx,
		middleware.SetupCORS(bannerHandler.GetBannersListHandler, configMux.addrOrigin, configMux.schema)))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(router, logger))

	return mux, nil
}
