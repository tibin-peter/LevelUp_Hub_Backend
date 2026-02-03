package user


type UpdateUserDTO struct{
	Name     string `json:"name"`
 	Email    string `json:"email"`
 	Password string `json:"password"` 
	ProfilePicURL string
}

// // Register dto
// type RegiserRequest struct {
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// //Login dto
// type LoginRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// //OTP varification dto
// type VerifyOTPRequest struct {
// 	Email string `jaon:"email"`
// 	Code  string `json:"code"`
// }
