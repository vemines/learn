package services

import (
	"gorm_tutorial/models"
	"log"

	"gorm.io/gorm"
)

func GetAllMessages(db *gorm.DB) []models.Message {
	var messages []models.Message
	err := db.Find(&messages).Error
	if err != nil {
		log.Fatalln("Cannot get all messages:", err)
	}

	return messages
}

func GetMessagesBySender(db *gorm.DB, userID uint) ([]models.Message, error) {
	var messages []models.Message

	// Execute the query.
	err := db.Model(&models.Message{}).
		Joins("INNER JOIN users ON messages.user_id_sender = users.user_id").
		Where("messages.user_id_sender = ?", userID).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessagesWithSenderAndRecipientNames gets all messages, along with the sender and recipient names.
func GetMessagesWithSenderAndRecipientNames(db *gorm.DB, senderID uint, recipientID uint) ([]models.MessageWithSenderAndRecipient, error) {
	var messages []models.MessageWithSenderAndRecipient

	// Execute the query.
	sql := `SELECT  m.message_id, m.content, m.message_status,
					u1.user_id AS sender_id, u1.username AS sender_name,
					u2.user_id AS recipient_id, u2.username AS recipient_name
			FROM messages AS m
			INNER JOIN users AS u1 ON m.user_id_sender = u1.user_id
			INNER JOIN users AS u2 ON m.user_id_recipient = u2.user_id
			WHERE (m.user_id_sender = ? AND m.user_id_recipient = ?)
			OR (m.user_id_sender = ? AND m.user_id_recipient = ?)`

	err := db.Raw(sql, senderID, recipientID, recipientID, senderID).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func GetMessagesWithSenderAndRecipientNamesFromView(db *gorm.DB) ([]models.MessageWithSenderAndRecipient, error) {
	var messages []models.MessageWithSenderAndRecipient

	// Execute the query.
	err := db.Raw("SELECT * FROM messages_with_sender_and_recipient_names").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
