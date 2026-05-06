package api

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"lbe_crypto_signer/internal/service"
	"lbe_crypto_signer/pkg/tools/log"
	"net"
	"net/http"
	"strconv"
)

func Start(ctx context.Context, apiConfig *service.ActiveConfig) error {

	log.ZInfo(ctx, "Start", zap.Any("ActiveConfig", apiConfig))
	router := NewGinRouter(apiConfig)
	address := net.JoinHostPort(apiConfig.ApiConf.Api.ListenIP, strconv.Itoa(apiConfig.ApiConf.Api.Ports))
	server := http.Server{Addr: address, Handler: router}

	log.ZInfo(ctx, "API server is initializing", zap.Any("ListenIP", apiConfig.ApiConf.Api.ListenIP), zap.Any("ApiPort", apiConfig.ApiConf.Api.Ports))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.ZInfo(ctx, "Start failed", zap.Any("error", err.Error()))
		netErr := errors.New(fmt.Sprintf("Api start err: %v", err))
		err = server.Shutdown(ctx)
		if err != nil {
			log.ZError(ctx, "Block Chain API Server shutdown failed.", err)
			return netErr
		}
	}
	return nil
}
