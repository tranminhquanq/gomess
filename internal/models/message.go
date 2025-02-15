package models

import "time"

type MessageType string
type ConversationType string
type AttachmentType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeVideo MessageType = "video"
	MessageTypeFile  MessageType = "file"

	ConversationTypeSingle ConversationType = "single"
	ConversationTypeGroup  ConversationType = "group"

	AttachmentTypeImage AttachmentType = "image"
	AttachmentTypeVideo AttachmentType = "video"
	AttachmentTypeFile  AttachmentType = "file"
)

type Message struct {
	ID             int64       `json:"id" db:"id"`
	ConversationID int64       `json:"conversation_id" db:"conversation_id"`
	SenderID       int64       `json:"sender_id" db:"sender_id"`
	Type           MessageType `json:"type" db:"type"`
	Message        string      `json:"message" db:"message"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at" db:"updated_at"`
}

func (u *Message) TableName() string {
	return "messages"
}

type Conversation struct {
	ID        int64            `json:"id" db:"id"`
	CreatorID int64            `json:"creator_id" db:"creator_id"`
	Title     string           `json:"title" db:"title"`
	Type      ConversationType `json:"type" db:"type"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at" db:"updated_at"`
}

func (c *Conversation) IsCreator(userID int64) bool {
	return c.CreatorID == userID
}

func (c *Conversation) IsSingle() bool {
	return c.Type == ConversationTypeSingle
}

func (c *Conversation) IsGroup() bool {
	return c.Type == ConversationTypeGroup
}

func (u *Conversation) TableName() string {
	return "conversations"
}

type Participant struct {
	ID             int64     `json:"id" db:"id"`
	ConversationID int64     `json:"conversation_id" db:"conversation_id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

func (u *Participant) TableName() string {
	return "participants"
}

type Attachment struct {
	ID        int64          `json:"id" db:"id"`
	MessageID int64          `json:"message_id" db:"message_id"`
	Type      AttachmentType `json:"type" db:"type"`
	URL       string         `json:"url" db:"url"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

func (u *Attachment) TableName() string {
	return "attachments"
}
