package adminprofile

type UpdateProfile struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type UpdateProfilePicture struct {
	ProfilePicURL string `json:"profile_Pic_url"`
}

type ChangePassword struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConformPassword string `json:"conform_password"`
}


type AdminProfileResponse struct {
	Name   string 
	Email  string 
	Role   string
	ProfilePicURL string
}

type AdminProfilePicResponse struct {
	ProfilePicURL string
}