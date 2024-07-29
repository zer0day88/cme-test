package route

import (
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/zer0day88/cme/message-service/internal/app/domain/repository"
	"github.com/zer0day88/cme/message-service/internal/app/handler"
	mw "github.com/zer0day88/cme/message-service/internal/app/middleware"
	"github.com/zer0day88/cme/message-service/internal/app/service"
)

func InitRoute(e *echo.Echo, log zerolog.Logger,
	cdb *gocql.Session, rdb *redis.Client) {

	userRepo := repository.NewUserRepository(cdb)
	msgRepo := repository.NewMessageRepository(cdb)
	msgService := service.NewMessageService(log, rdb, *msgRepo, *userRepo)
	msgHandler := handler.NewUserHandler(log, msgService)
	mAuth := mw.InitMAuth(rdb)

	e.Use(middleware.CORS())

	r := e.Group("/v1")
	r.GET("/ok", msgHandler.Ok)
	r.Use(mAuth.JWT)
	{
		r.POST("/send", msgHandler.Send)
		r.GET("/messages", msgHandler.GetMessageHistory)
	}
}
