package message

import "time"

//for get the targeruserid
type StartConversationRequest struct{
	TargetUserID  uint `json:"target_user_id"`
}

type SendWSMessage struct {
	Content string `json:"content"`
}

//for list the conversation sidebar
type ConversationWithLastMsg struct {
	ConversationID uint
	OtherUserID    uint
	OtherUserName  string    
	OtherProfilePic string
	LastMessage    string
	LastTime       *time.Time
}