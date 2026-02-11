package message

import "gorm.io/gorm"

type Repository interface {
	FindByUsers(u1,u2 uint)(*Conversation,error)
	CreateConversation(c *Conversation)error
	GetConversationByID(id uint)(*Conversation,error)
	CreateMessage(m *Message) error
	ListUserConversations(userID uint) ([]ConversationWithLastMsg, error)
	GetMessages(cid uint, beforeID uint) ([]Message, error)
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) FindByUsers(u1,u2 uint)(*Conversation,error){
	var c Conversation

	err := r.db.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		u1,u2,u2,u1,
	).First(&c).Error

	if err == gorm.ErrRecordNotFound {
		return nil,nil
	}

	return &c,err
}

func (r *repo) CreateConversation(c *Conversation)error{
	return r.db.Create(c).Error
}

func (r *repo) GetConversationByID(id uint)(*Conversation,error){
	var c Conversation
	err:=r.db.First(&c,id).Error
	if err==gorm.ErrRecordNotFound{
		return nil,nil
	}
	return &c,err
}

func (r *repo) CreateMessage(m *Message) error {
	return r.db.Create(m).Error
}

func (r *repo) ListUserConversations(userID uint) ([]ConversationWithLastMsg, error) {

	var result []ConversationWithLastMsg

	err := r.db.Raw(`
	SELECT 
		c.id as conversation_id,
		CASE 
			WHEN c.user1_id = ? THEN c.user2_id
			ELSE c.user1_id
		END as other_user_id,
		m.content as last_message,
		m.created_at as last_time
	FROM conversations c
	LEFT JOIN messages m 
		ON m.id = (
			SELECT id FROM messages 
			WHERE conversation_id = c.id
			ORDER BY created_at DESC
			LIMIT 1
		)
	WHERE c.user1_id = ? OR c.user2_id = ?
	ORDER BY m.created_at DESC
	`, userID, userID, userID).
	Scan(&result).Error

	return result, err
}

//for get message for history
func (r *repo) GetMessages(cid uint, beforeID uint) ([]Message, error) {

	var msgs []Message

	q := r.db.
		Where("conversation_id = ?", cid)

	if beforeID != 0 {
		q = q.Where("id < ?", beforeID)
	}

	err := q.
		Order("id ASC").
		Limit(20).
		Find(&msgs).Error

	return msgs, err
}