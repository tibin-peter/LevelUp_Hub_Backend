package user

import "fmt"

type Service interface {
	GetUserById(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// func for findbyid
func (s *service) GetUserById(id uint) (*User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	user,err:=s.repo.FindById(id)
	if err!=nil{
		return nil,fmt.Errorf("user not found")
	}
	return user,nil
}

// func for findbyemail
func (s *service) FindByEmail(email string) (*User, error) {
	if email==""{
		return nil,fmt.Errorf("email is required")
	}
		user,err:=s.repo.FindByEmail(email)
		if err!=nil{
			return nil,fmt.Errorf("user not found")
		}
	return user,nil
}

// for update
func (s *service) Update(user *User) error {
	if user.ID==0{
		return fmt.Errorf("invalid user id")
	}
	return s.repo.Update(user)
}

// for delete
func (s *service) Delete(id uint) error {
	if id==0{
		return fmt.Errorf("invalid user id")
	}
	if err:=s.repo.Delete(id);err!=nil{
		return fmt.Errorf("user not found")
	}
	return nil
}