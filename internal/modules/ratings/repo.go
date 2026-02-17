package ratings

import "gorm.io/gorm"

type Repository interface {
	Create(rating *Rating) error
	ExistsByBooking(bookingID uint) (bool, error)
	GetByMentor(mentorID uint) ([]Rating, error)
	GetAverageByMentor(mentorID uint) (float64, error)
	CountByMentor(mentorID uint) (int64, error)
	GetTopMentors() ([]TopMentor, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

//create
func (r *repo) Create(rating *Rating) error {
	return r.db.Create(rating).Error
}

//check exist prevent duplicate
func (r *repo) ExistsByBooking(bookingID uint) (bool, error) {
	var count int64

	err := r.db.Model(&Rating{}).
		Where("booking_id = ?", bookingID).
		Count(&count).Error

	return count > 0, err
}

//get ratings of mentor
func (r *repo) GetByMentor(mentorProfileID uint) ([]Rating, error) {
	var ratings []Rating

	err := r.db.
		Where("mentor_profile_id = ?", mentorProfileID).
		Order("created_at DESC").
		Find(&ratings).Error

	return ratings, err
}

//get mentor avarage rating
func (r *repo) GetAverageByMentor(mentorProfileID uint) (float64, error) {
	var avg float64

	err := r.db.Model(&Rating{}).
		Where("mentor_profile_id = ?", mentorProfileID).
		Select("AVG(rating)").
		Scan(&avg).Error

	return avg, err
}

//for count
func (r *repo) CountByMentor(mentorProfileID uint) (int64, error) {
	var count int64

	err := r.db.Model(&Rating{}).
		Where("mentor_profile_id = ?", mentorProfileID).
		Count(&count).Error

	return count, err
}

//for get top mentors
func (r *repo) GetTopMentors() ([]TopMentor, error) {
	var result []TopMentor

	err := r.db.Table("ratings").
		Select("mentor_profile_id, AVG(rating) as avg_rating").
		Group("mentor_profile_id").
		Having("AVG(rating) >= ?", 4).
		Order("avg_rating DESC").
		Scan(&result).Error

	return result, err
}
