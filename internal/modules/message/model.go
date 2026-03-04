package message

import "time"

type Conversation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User1ID   uint      `gorm:"uniqueIndex:idx_users" json:"user1_id"`
	User2ID   uint      `gorm:"uniqueIndex:idx_users" json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID uint      `gorm:"index;not null" json:"conversation_id"`
	SenderID       uint      `gorm:"index;not null" json:"sender_id"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
