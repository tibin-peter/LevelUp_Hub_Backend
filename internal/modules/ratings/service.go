package ratings

import (
	"LevelUp_Hub_Backend/internal/modules/booking"
	"errors"
)

type Service interface {
	CreateRating(studentID uint,req CreateRatingRequest)(*CreateRatingResponse, error)
	GetMentorRatings(mentorID uint) ([]RatingResponse, error)
	GetMentorSummary(mentorID uint) (*MentorRatingSummary, error)
	GetTopMentors() ([]TopMentor, error)
}

type service struct {
	repo         Repository
	bookingrepo  booking.Repository
}

func NewService(r Repository,b booking.Repository)Service{
	return &service{
		repo: r,
		bookingrepo: b,
	}
}

//create ratings
func(s *service)CreateRating(studentID uint,req CreateRatingRequest)(*CreateRatingResponse, error){
	// Get booking
	booking, err := s.bookingrepo.GetByID(req.BookingID)
	if err != nil {
		return nil,errors.New("booking not found")
	}

	//Check ownership
	if booking.StudentID != studentID {
		return nil,errors.New("not your booking")
	}

	// Check completed
	if booking.Status != "completed" {
		return nil,errors.New("session not completed yet")
	}

	//Check duplicate rating
	exists, err := s.repo.ExistsByBooking(req.BookingID)
	if err != nil {
		return nil,err
	}
	if exists {
		return nil,errors.New("rating already submitted")
	}

	//Validate rating value
	if req.Rating < 1 || req.Rating > 5 {
		return nil,errors.New("rating must be between 1 and 5")
	}

	//Create rating
	rating := &Rating{
		BookingID: req.BookingID,
		StudentID: studentID,
		MentorProfileID:  booking.MentorProfileID,
		Rating:    req.Rating,
		Review:    req.Review,
	}

	if err := s.repo.Create(rating); err != nil {
		return nil, err
	}
	avg, err := s.repo.GetAverageByMentor(booking.MentorProfileID)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountByMentor(booking.MentorProfileID)
	if err != nil {
		return nil, err
	}
	res:=&CreateRatingResponse{
		AvgRating: avg,
		TotalReviews: count,
	}
	return res,nil
}

//get mentor ratings
func(s *service)GetMentorRatings(mentorID uint)([]RatingResponse,error){
	ratings, err := s.repo.GetByMentor(mentorID)
	if err != nil {
		return nil, err
	}

	var res []RatingResponse
	for _, r := range ratings {
		res = append(res, RatingResponse{
			StudentID: r.StudentID,
			Rating:    r.Rating,
			Review:    r.Review,
			CreatedAt: r.CreatedAt,
		})
	}

	return res, nil
}

//get mnetor avg
func(s *service)GetMentorSummary(mentorID uint)(*MentorRatingSummary,error){
	avg, err := s.repo.GetAverageByMentor(mentorID)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountByMentor(mentorID)
	if err != nil {
		return nil, err
	}

	return &MentorRatingSummary{
		AvgRating:    avg,
		TotalReviews: count,
	}, nil
}

//get top mentors
func (s *service)GetTopMentors()([]TopMentor, error){
	return s.repo.GetTopMentors()
}
