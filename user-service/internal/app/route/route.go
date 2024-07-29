package route

import (
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/zer0day88/cme/user-service/internal/app/domain/repository"
	"github.com/zer0day88/cme/user-service/internal/app/handler"
	"github.com/zer0day88/cme/user-service/internal/app/service"
)

func InitRoute(e *echo.Echo, log zerolog.Logger,
	cdb *gocql.Session, rdb *redis.Client) {

	userRepo := repository.NewUserRepository(cdb)
	userService := service.NewUserService(log, rdb, *userRepo)
	userHandler := handler.NewUserHandler(log, userService)

	e.Use(middleware.CORS())

	r := e.Group("/v1")
	{
		r.GET("/ok", userHandler.Ok)
		r.POST("/register", userHandler.Register)
		r.POST("/login", userHandler.Login)
	}
}
