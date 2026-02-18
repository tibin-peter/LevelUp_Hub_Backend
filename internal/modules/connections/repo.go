package connections

import "gorm.io/gorm"

type Repository interface {
	Create(c *Connection) error
	Exists(studentID, mentorProfileID uint) (bool, error)

	CountStudents(mentorProfileID uint) (int64, error)
	GetConnectedMentors(studentID uint) ([]ConnectedMentor, error)
	GetConnectedStudents(mentorProfileID uint,) ([]ConnectedStudent, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(c *Connection) error {
	return r.db.Create(c).Error
}

func (r *repo) Exists(studentID, mentorProfileID uint) (bool, error) {
	var count int64

	err := r.db.Model(&Connection{}).
		Where("student_id=? AND mentor_profile_id=?", studentID, mentorProfileID).
		Count(&count).Error

	return count > 0, err
}

func (r *repo) CountStudents(mentorProfileID uint) (int64, error) {
	var count int64

	err := r.db.Model(&Connection{}).
		Where("mentor_profile_id=?", mentorProfileID).
		Count(&count).Error

	return count, err
}

func (r *repo) GetConnectedMentors(studentID uint) ([]ConnectedMentor, error) {

	var res []ConnectedMentor

	err := r.db.Table("connections c").
		Select(`
			mp.id as mentor_profile_id,
			u.name,
			mp.profile_pic_url,
			mp.category,
			mp.hourly_price,
			mp.rating_avg
		`).
		Joins("JOIN mentor_profiles mp ON mp.id = c.mentor_profile_id").
		Joins("JOIN users u ON u.id = mp.user_id").
		Where("c.student_id = ?", studentID).
		Scan(&res).Error

	return res, err
}

func (r *repo) GetConnectedStudents(mentorProfileID uint,) ([]ConnectedStudent, error) {

	var res []ConnectedStudent

	err := r.db.Table("connections c").
		Select(`
			u.id as student_id,
			u.name,
			u.profile_pic_url
		`).
		Joins("JOIN users u ON u.id = c.student_id").
		Where("c.mentor_profile_id = ?", mentorProfileID).
		Scan(&res).Error

	return res, err
}