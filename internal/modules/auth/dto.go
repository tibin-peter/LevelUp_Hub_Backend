package auth


// Register dto
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	OTP  string `json:"otp"`
}

//Login dto
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Send otp
type SendOTPRequest struct{
	Email string `json:"email"`
}

//for register response
type AuthResponseData struct {
    Email string `json:"email"`
    Role  string `json:"role"`
    Name  string `json:"name"`
		IsVerified bool
		ProfilePicURL string
}