package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/zer0day88/cme/message-service/helper"
	"github.com/zer0day88/cme/message-service/internal/app/domain/entities"
	"github.com/zer0day88/cme/message-service/internal/app/domain/repository"
	"github.com/zer0day88/cme/message-service/internal/app/model"
	"github.com/zer0day88/cme/message-service/pkg/response"
	"time"
)

type MessageService struct {
	log      zerolog.Logger
	msgRepo  repository.MessageRepository
	userRepo repository.UserRepository
	rdb      *redis.Client
}

func NewMessageService(log zerolog.Logger, rdb *redis.Client, msgRepo repository.MessageRepository, userRepo repository.UserRepository) *MessageService {
	return &MessageService{
		log:      log,
		rdb:      rdb,
		msgRepo:  msgRepo,
		userRepo: userRepo,
	}
}

func (s *MessageService) GetByRecipient(ctx context.Context, userID string) ([]entities.Message, response.ApiJSON) {

	user, err := s.userRepo.FindOneByID(ctx, userID)

	if err != nil {
		s.log.Err(err).Send()
		return nil, response.ErrSystem
	}

	msg := make([]entities.Message, 0)
	err = s.msgRepo.GetByRecipient(ctx, user.Username, &msg)

	if err != nil {
		s.log.Err(err).Send()
		return nil, response.ErrSystem.WithErr(err)
	}
	return msg, response.OKNoError
}

func (s *MessageService) Send(ctx context.Context, userID string, send model.SendRequest) response.ApiJSON {
	errv := helper.ValidateStruct(send)

	if errv != nil {
		return response.ErrBadRequest.WithMsg(*errv)
	}

	user, err := s.userRepo.FindOneByID(ctx, userID)

	if err != nil {

		s.log.Err(err).Send()
		return response.ErrSystem
	}

	msgID := uuid.NewString()

	err = s.msgRepo.Insert(ctx,
		entities.Message{
			ID: msgID, Sender: user.Username, Recipient: send.To, Content: send.Content, MessageTime: time.Now(),
		})

	if err != nil {
		s.log.Err(err).Send()
		return response.ErrSystem.WithErr(err)
	}
	return response.OKNoError
}
