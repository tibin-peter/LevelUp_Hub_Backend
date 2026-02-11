package message

import "time"

type Conversation struct {
	ID        uint `gorm:"primaryKey"`
	User1ID uint `gorm:"uniqueIndex:idx_users"`
	User2ID uint `gorm:"uniqueIndex:idx_users"`
	CreatedAt time.Time
}

type Message struct {
	ID             uint   `gorm:"primaryKey"`
	ConversationID uint   `gorm:"index;not null"`
	SenderID       uint   `gorm:"index;not null"`
	Content        string `gorm:"type:text;not null"`
	CreatedAt      time.Time
}
