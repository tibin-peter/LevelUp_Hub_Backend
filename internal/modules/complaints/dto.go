package complaints


type CreateComplaintRequest struct {
	Category    string `json:"category"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type ReplyRequest struct {
	Reply  string `json:"reply"`
	Status string `json:"status"`
}