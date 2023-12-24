package models

import "time"

type Message struct {
	MessageId         uint          `gorm:"primary_key;column:message_id"`
	Text              string        `gorm:"column:content"`
	User_Id_Sender    int           `gorm:"column:user_id_sender"`
	User_Id_Recipient int           `gorm:"column:user_id_recipient"`
	CreatedAt         time.Time     `gorm:"column:created_at"`
	LastUpdatedAt     time.Time     `gorm:"column:last_updated_at"`
	Status            MessageStatus `gorm:"column:message_status"`
}

type MessageWithSenderAndRecipient struct {
	MessageID     uint          `json:"message_id"`
	Content       string        `json:"content"`
	SenderID      uint          `json:"sender_id"`
	SenderName    string        `json:"sender_name"`
	RecipientID   uint          `json:"recipient_id"`
	RecipientName string        `json:"recipient_name"`
	Status        MessageStatus `json:"message_status"`
}

type MessageStatus string // enum represents the status of a message.

const (
	MessageStatusUnread  MessageStatus = "Unread"
	MessageStatusRead    MessageStatus = "Read"
	MessageStatusDeleted MessageStatus = "Deleted"
)
