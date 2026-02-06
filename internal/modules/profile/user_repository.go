package profile

import (
	"LevelUp_Hub_Backend/internal/repository/generic"
	"errors"

	"gorm.io/gorm"
)


type Repository interface {
	CreateUser(user *User)error
	FindUserById(id uint)(*User,error)
	FindUserByEmail(email string)(*User,error)
	UpdateUser(user *User)error
	DeleteUser(id uint)error
}

type repo struct {
	base *generic.Repository[User]
	db *gorm.DB
}

func NewRepository(db *gorm.DB)Repository{
	return &repo{
		base: generic.NewRepository[User](db),
		db: db,
	}
}

//create user
func (r *repo)CreateUser(user *User)error{
	return r.base.Create(user)
}
//find user by id
func (r *repo)FindUserById(id uint)(*User,error){
	return r.base.FindById(id)
}
//for user Update
func (r *repo)UpdateUser(user *User)error{
	return r.base.Update(user)
}
//for delete user
func (r *repo)DeleteUser(id uint)error{
	return r.base.Delete(id)
}

//for specific Query
func (r *repo)FindUserByEmail(email string)(*User,error){
	var user User
	err:=r.db.Where("email = ?",email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, nil // no user, no error
	}
	return &user,err
}