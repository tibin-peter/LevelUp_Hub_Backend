package dashboard

import (
	"LevelUp_Hub_Backend/internal/modules/booking"
	"LevelUp_Hub_Backend/internal/modules/connections"
	"LevelUp_Hub_Backend/internal/modules/favorites"
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/ratings"
)

type Service interface {
	StudentDashboard(userID uint) (*StudentDashboard, error)
	MentorDashboard(profileID uint) (*MentorDashboard,error)
}

type service struct {
	bookingRepo booking.Repository
	favoritesRepo favorites.Repository
	ratingRepo  ratings.Repository
	paymentRepo payment.Repository
	connectionRepo connections.Repository
	mentorRepo profile.MentorRepository
}

func NewService(
	b booking.Repository,
	fav favorites.Repository,
	r  ratings.Repository,
	pay payment.Repository,
	c connections.Repository,
	m profile.MentorRepository,
) Service {
	return &service{
		bookingRepo: b,
		favoritesRepo: fav,
		ratingRepo: r,
		paymentRepo: pay,
		connectionRepo: c,
		mentorRepo: m,
	}
}

//////////////////// Student Dashboard  //////////////////
func (s *service) StudentDashboard(userID uint) (*StudentDashboard, error) {

	active,_ := s.bookingRepo.CountByStatus(userID,"confirmed")
	completed,_ := s.bookingRepo.CountByStatus(userID,"completed")
	favs,_ := s.favoritesRepo.CountByStudent(userID)

	upcoming,_ := s.bookingRepo.GetUpcomingByStudent(userID)
	var upcomingDTO []UpcomingSession

for _, u := range upcoming {
	upcomingDTO = append(upcomingDTO, UpcomingSession{
		BookingID:   u.BookingID,
		MentorName:  u.MentorName,
		StartTime: u.StartTime,
		EndTime: u.EndTime,
	})
}
	mentors,_ := s.connectionRepo.GetConnectedMentors(userID)
	var mentorDTO []ConnectedMentor

for _, m := range mentors {
	mentorDTO = append(mentorDTO, ConnectedMentor{
		MentorProfileID: m.MentorProfileID,
		Name:            m.Name,
		ProfilePicURL:   m.ProfilePicURL,
		Category:        m.Category,
		HourlyPrice:     m.HourlyPrice,
		RatingAvg:       m.RatingAvg,
	})
}

	return &StudentDashboard{
		ActiveBookings: active,
		CompletedSessions: completed,
		FavoriteMentors: favs,
		UpcomingSessions: upcomingDTO,
		ConnectedMentors: mentorDTO,
	},nil
}

//////////////////// Mentor Dashboard  //////////////////
func (s *service) MentorDashboard(userID uint) (*MentorDashboard,error){

	profile, err := s.mentorRepo.FindMentorByUserID(userID)
	if err != nil {
		return nil, err
	}

	profileID := profile.ID

	earnings,_ := s.paymentRepo.SumByMentor(profileID)
	students,_ := s.connectionRepo.CountStudents(profileID)
	rating,_ := s.ratingRepo.GetAverageByMentor(profileID)
	requests,_ := s.bookingRepo.CountRequests(profileID)

	return &MentorDashboard{
		TotalEarnings: earnings,
		TotalStudents: students,
		AvgRating: rating,
		BookingRequests: requests,
	},nil
}