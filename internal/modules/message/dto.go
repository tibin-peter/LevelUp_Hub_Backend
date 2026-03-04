package message

import "time"

// StartConversationRequest for initiating a new chat
type StartConversationRequest struct {
	TargetUserID uint `json:"target_user_id"`
}

// SendWSMessage for messages sent via WebSocket
type SendWSMessage struct {
	Content string `json:"content"`
}

// ConversationWithLastMsg for the sidebar inbox list
type ConversationWithLastMsg struct {
	ConversationID  uint       `json:"conversation_id"`
	OtherUserID      uint       `json:"other_user_id"`
	OtherUserName    string     `json:"other_user_name"`
	OtherProfilePic string     `json:"other_profile_pic"`
	LastMessage     string     `json:"last_message"`
	LastTime        *time.Time `json:"last_time"`
}