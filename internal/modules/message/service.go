package message

type Service interface{
	StartConversation(uid,tid uint)(uint,error)
	SendMessage(cid,sender uint,content string) error
	IsUserInConversation(cid,uid uint)(bool,error)
	GetOtherUser(cid,uid uint)(uint,error)
	ListUserConversations(uid uint) ([]ConversationWithLastMsg, error)
	GetMessages(cid uint, before uint) ([]Message, error)
}

type service struct{ repo Repository }

func NewService(r Repository)Service{
	return &service{r}
}

func (s *service) StartConversation(uid,tid uint)(uint,error){
	u1,u2 := uid,tid
	if u1>u2 { u1,u2=u2,u1 }

	exist,_ := s.repo.FindByUsers(u1,u2)
	if exist!=nil {
		return exist.ID,nil
	}

	c := Conversation{
		User1ID:u1,
		User2ID:u2,
	}
	if err:=s.repo.CreateConversation(&c);err!=nil{
		return 0,err
	}
	return c.ID,nil
}

func (s *service) SendMessage(cid,sender uint,content string) error{
	return s.repo.CreateMessage(&Message{
		ConversationID:cid,
		SenderID:sender,
		Content:content,
	})
}

func (s *service) IsUserInConversation(cid,uid uint)(bool,error){
	c,err := s.repo.GetConversationByID(cid)
	if err!=nil||c==nil { return false,err }
	return c.User1ID==uid||c.User2ID==uid,nil
}

func (s *service) GetOtherUser(cid,uid uint)(uint,error){
	c,_ := s.repo.GetConversationByID(cid)
	if c.User1ID==uid { return c.User2ID,nil }
	return c.User1ID,nil
}

func (s *service) ListUserConversations(uid uint) ([]ConversationWithLastMsg, error) {
	return s.repo.ListUserConversations(uid)
}

func (s *service) GetMessages(cid uint, before uint) ([]Message, error) {
	return s.repo.GetMessages(cid, before)
}