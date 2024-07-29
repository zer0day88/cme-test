package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zer0day88/cme/user-service/config"
	"github.com/zer0day88/cme/user-service/infra/db"
	"github.com/zer0day88/cme/user-service/internal/app/route"
	"github.com/zer0day88/cme/user-service/pkg/logger"
	"github.com/zer0day88/cme/user-service/pkg/shutdown"
	"time"
)

func main() {

	config.Load()

	log := logger.New()

	cdb, err := db.InitCassandra()
	if err != nil {
		panic(err)
	}

	rdb, err := db.InitRedis()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	route.InitRoute(e, log, cdb, rdb)

	fmt.Println(config.Key.ServerPort)

	go func() {
		if err := e.Start(":" + config.Key.ServerPort); err != nil {
			log.Fatal().Msgf("shutting down: %s", err)
		}
	}()

	wait := shutdown.GracefulShutdown(context.Background(), log, 10*time.Second, map[string]shutdown.Operation{
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	<-wait
}
