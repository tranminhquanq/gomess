package factory

import (
	"time"

	"github.com/tranminhquanq/gomess/internal/app/domain"
)

type MessageFactory struct{}

func (m MessageFactory) CreateMessage(
	id int64,
	conversationId int64,
	senderId int64,
	message string,
) domain.Message {
	return domain.Message{
		ID:             id,
		ConversationID: conversationId,
		SenderID:       senderId,
		Message:        message,
		CreatedAt:      time.Now(),
	}
}
