package repository

import "github.com/tranminhquanq/gomess/internal/app/domain"

type MessageRepository interface {
	SaveMessage(domain.Message) (domain.Message, error)
	FindMessagesInConversation(conversationId int64, offset, limit int) (domain.ListResult[domain.Message], error)
}
