package usermanagement

type PendingMentorDTO struct {
	MentorProfileID uint
	UserID          uint
	Name            string
	Email           string
	Status          string
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}