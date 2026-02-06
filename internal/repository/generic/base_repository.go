package generic

import (
	"errors"

	"gorm.io/gorm"
)

//for reuse basic crud
type Repository[T any] struct {
	db *gorm.DB
}

//dep injection
func NewRepository[T any](db *gorm.DB)*Repository[T]{
	return &Repository[T]{db}
}

//create
func (r *Repository[T])Create(entity *T)error{
	return r.db.Create(entity).Error
}

//finc by id
func (r *Repository[T]) FindById(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &entity, err
}

//for findall
func (r *Repository[T])FindAll()([]T,error){
	var list []T
	err:=r.db.Find(&list).Error
	return list,err
}

//update
func (r *Repository[T])Update(entity *T)error{
	return r.db.Save(entity).Error
}

//delete
func (r *Repository[T])Delete(id uint)error{
	var entity T
	return r.db.Delete(&entity,id).Error
}