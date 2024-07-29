package service

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/zer0day88/cme/user-service/config"
	"github.com/zer0day88/cme/user-service/helper"
	"github.com/zer0day88/cme/user-service/internal/app/domain/entities"
	"github.com/zer0day88/cme/user-service/internal/app/domain/repository"
	"github.com/zer0day88/cme/user-service/internal/app/model"
	"github.com/zer0day88/cme/user-service/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	log      zerolog.Logger
	userRepo repository.UserRepository
	rdb      *redis.Client
}

func NewUserService(log zerolog.Logger, rdb *redis.Client, userRepo repository.UserRepository) *UserService {
	return &UserService{
		log:      log,
		rdb:      rdb,
		userRepo: userRepo,
	}
}

func (s *UserService) Get(ctx context.Context) ([]entities.User, response.ApiJSON) {

	user := make([]entities.User, 0)
	err := s.userRepo.Get(ctx, &user)

	if err != nil {
		s.log.Err(err).Send()
		return nil, response.ErrSystem.WithErr(err)
	}
	return user, response.OKNoError
}

func (s *UserService) Register(ctx context.Context, register model.RegisterRequest) response.ApiJSON {
	errv := helper.ValidateStruct(register)

	if errv != nil {
		return response.ErrBadRequest.WithMsg(*errv)
	}

	//if ok, errs := helper.ValidatePassword(register.Password); !ok {
	//	return response.ErrBadRequest.WithMsg(errs)
	//}

	row, err := s.userRepo.FindOneByUsername(ctx, register.Username)

	if err != nil && !errors.Is(err, gocql.ErrNotFound) {

		s.log.Err(err).Send()
		return response.ErrSystem
	}

	if row != nil {
		return response.ErrBadRequest.WithMsg("Username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), 8)
	if err != nil {
		s.log.Err(err).Send()
		return response.ErrSystem.WithErr(err)
	}

	id := uuid.NewString()

	err = s.userRepo.Insert(ctx, entities.User{ID: id, Username: register.Username, Password: string(hashedPassword), CreatedAt: time.Now()})

	if err != nil {
		s.log.Err(err).Send()
		return response.ErrSystem.WithErr(err)
	}
	return response.OKNoError
}

func (s *UserService) Login(ctx context.Context, login model.LoginRequest) (*helper.JwtToken, response.ApiJSON) {
	//validate
	errv := helper.ValidateStruct(login)

	if errv != nil {
		return nil, response.ErrBadRequest.WithMsg(*errv)
	}

	authData, err := s.userRepo.FindOneByUsername(ctx, login.Username)

	if err != nil {
		return nil, response.ErrUnauthorized.WithErr(err)
	}

	//compare password
	err = bcrypt.CompareHashAndPassword([]byte(authData.Password), []byte(login.Password))
	if err != nil {

		//TODO: increase login failed attempt
		//if reached limit, in this case 5, block until specific time

		return nil, response.ErrUnauthorized.WithErr(err)
	}

	var (
		accessTokenExpiresIn   = time.Minute * config.Key.AccessTokenExpiredIn
		accessTokenPrivateKey  = config.Key.AccessTokenPrivateKey
		refreshTokenExpiresIn  = time.Minute * config.Key.RefreshTokenExpiredIn
		refreshTokenPrivateKey = config.Key.RefreshTokenPrivateKey
	)

	//create Token
	accessTokenDetail, err := helper.CreateToken(authData.ID, accessTokenExpiresIn, accessTokenPrivateKey)
	if err != nil {
		return nil, response.ErrSystem.WithErr(err)
	}

	refreshTokenDetail, err := helper.CreateToken(authData.ID, refreshTokenExpiresIn, refreshTokenPrivateKey)
	if err != nil {
		return nil, response.ErrSystem.WithErr(err)
	}

	jwtToken := helper.JwtToken{
		AccessToken:  accessTokenDetail,
		RefreshToken: refreshTokenDetail,
	}

	// for none blocking task such as send email, update last login, etc...
	// or passing the task to event streaming using kafka, NATS or pub/sub
	go func() {

		ctxBg := context.Background()

		jwtKeyAccess := fmt.Sprintf("%s:%s", config.Key.JwtRedisKey, accessTokenDetail.TokenUuid)
		jwtKeyRefresh := fmt.Sprintf("%s:%s", config.Key.JwtRedisKey, refreshTokenDetail.TokenUuid)

		//Save token data
		err = s.rdb.Set(ctxBg, jwtKeyAccess, authData.ID, config.Key.AccessTokenExpiredIn*time.Minute).Err()

		if err != nil {
			s.log.Err(err).Send()
		}

		err = s.rdb.Set(ctxBg, jwtKeyRefresh, authData.ID, config.Key.RefreshTokenExpiredIn*time.Minute).Err()
		if err != nil {
			s.log.Err(err).Send()
		}

	}()

	return &jwtToken, response.OKNoError
}
