package repository

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/zer0day88/cme/message-service/internal/app/domain/entities"
)

type MessageRepository struct {
	db *gocql.Session
}

func NewMessageRepository(db *gocql.Session) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) GetByRecipient(ctx context.Context, recipient string, msg *[]entities.Message) error {
	scanner := r.db.Query(
		"select id,sender,recipient,content,message_time from messages where recipient = ?",
		recipient).
		WithContext(ctx).Iter().Scanner()

	for scanner.Next() {
		var m entities.Message
		err := scanner.Scan(&m.ID, &m.Sender, &m.Recipient, &m.Content, &m.MessageTime)
		if err != nil {

			return err
		}
		*msg = append(*msg, m)
	}
	if err := scanner.Err(); err != nil {

		return err
	}

	return nil
}

func (r *MessageRepository) Insert(ctx context.Context, user entities.Message) error {
	err := r.db.Query("insert into messages(id,sender,recipient,content,message_time) values (?,?,?,?,?)",
		user.ID, user.Sender, user.Recipient, user.Content, user.MessageTime).
		WithContext(ctx).Exec()

	if err != nil {
		return err
	}

	return nil
}
