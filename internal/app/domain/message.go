package domain

import (
	"time"

	"github.com/tranminhquanq/gomess/internal/models"
)

type Message struct {
	ID             int64              `json:"id"`
	ConversationID int64              `json:"conversation_id"`
	SenderID       int64              `json:"sender_id"`
	Type           models.MessageType `json:"type"`
	Message        string             `json:"message"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at,omitempty"`

	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	ID        int64                 `json:"id"`
	Type      models.AttachmentType `json:"type"`
	URL       string                `json:"url"`
	CreatedAt time.Time             `json:"created_at"`
}
