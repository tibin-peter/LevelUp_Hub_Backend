package user

import (
	"LevelUp_Hub_Backend/internal/repository/generic"

	"gorm.io/gorm"
)


type Repository interface {
	Create(user *User)error
	FindById(id uint)(*User,error)
	FindByEmail(email string)(*User,error)
	Update(user *User)error
	Delete(id uint)error
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
func (r *repo)Create(user *User)error{
	return r.base.Create(user)
}
//find user by id
func (r *repo)FindById(id uint)(*User,error){
	return r.base.FindById(id)
}
//for user Update
func (r *repo)Update(user *User)error{
	return r.base.Update(user)
}
//for delete user
func (r *repo)Delete(id uint)error{
	return r.base.Delete(id)
}

//for specific Query
func (r *repo)FindByEmail(email string)(*User,error){
	var user User
	err:=r.db.Where("email = ?",email).First(&user).Error
	return &user,err
}