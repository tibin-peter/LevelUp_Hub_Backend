package favorites

import "gorm.io/gorm"

type Repository interface {
	Create(f *Favorite) error
	Delete(studentID, mentorProfileID uint) error
	Exists(studentID, mentorProfileID uint) (bool, error)
	GetDetailedByStudent(studentID uint) ([]FavoriteResponse, error)
	CountByStudent(userID uint) (int64, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

//create
func (r *repo) Create(f *Favorite) error {
	return r.db.Create(f).Error
}

//delete
func (r *repo) Delete(studentID, mentorProfileID uint) error {
	return r.db.
		Where("student_id = ? AND mentor_profile_id = ?", studentID, mentorProfileID).
		Delete(&Favorite{}).Error
}

//check exist
func (r *repo) Exists(studentID, mentorProfileID uint) (bool, error) {
	var count int64

	err := r.db.Model(&Favorite{}).
		Where("student_id = ? AND mentor_profile_id = ?", studentID, mentorProfileID).
		Count(&count).Error

	return count > 0, err
}

//list favorite
func (r *repo) GetDetailedByStudent(studentID uint) ([]FavoriteResponse, error) {

	var res []FavoriteResponse

	err := r.db.Table("favorites f").
		Select(`
			mp.id as mentor_profile_id,
			u.name,
			mp.profile_pic_url,
			mp.category,
			mp.hourly_price,
			mp.rating_avg
		`).
		Joins("JOIN mentor_profiles mp ON mp.id = f.mentor_profile_id").
		Joins("JOIN users u ON u.id = mp.user_id").
		Where("f.student_id = ?", studentID).
		Scan(&res).Error

	return res, err
}

func (r *repo) CountByStudent(userID uint) (int64, error) {

	var count int64

	err := r.db.Model(&Favorite{}).
		Where("student_id=?", userID).
		Count(&count).Error

	return count, err
}