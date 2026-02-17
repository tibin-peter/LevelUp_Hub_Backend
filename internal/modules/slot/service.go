package slot

import (
	"errors"
	"time"
)

type Service interface {
	CreateSlot(userID uint, req CreateSlotRequest) error
	GetMentorSlots(userID uint) ([]MentorSlot, error)
	GetAvailableSlots(mentoID uint, date *time.Time) ([]MentorSlot, error)
	DeleteSlot(userID uint, slotID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// create slot
func (s *service) CreateSlot(userID uint, req CreateSlotRequest) error {
	now := time.Now()

	// 1. Lead time check: Must be at least 1 hour in the future
	if req.StartTime.Before(now.Add(time.Hour)) {
		return errors.New("slots must be created at least 1 hour in advance")
	}

	// 2. Duration check: Must be at least 1 hour
	if req.EndTime.Sub(req.StartTime) < time.Hour {
		return errors.New("slot duration must be at least 1 hour")
	}

	// 3. Basic validity check
	if !req.EndTime.After(req.StartTime) {
		return errors.New("end time must be after start time")
	}
	profileID, err := s.repo.GetProfileIDByUserID(userID)
	if err != nil {
		return errors.New("mentor profile not found")
	}
	//overlap check
	overlap, err := s.repo.HasOverlap(profileID, req.StartTime, req.EndTime)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("slot already existing in the current time")
	}
	//create slot
	slot := MentorSlot{
		MentorProfileID: profileID,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		IsBooked:        false,
	}
	return s.repo.CreateSlot(&slot)
}

// for mentor see theirown slot
func (s *service) GetMentorSlots(userID uint) ([]MentorSlot, error) {
	profileID, err := s.repo.GetProfileIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSlotsByMentor(profileID)
}

// for student see available slot
func (s *service) GetAvailableSlots(mentorID uint, date *time.Time) ([]MentorSlot, error) {
	if date == nil {
		return s.repo.GetAvailableSlotsByMentor(mentorID)
	}

	//filter by dare range
	y, m, d := date.Date()
	loc := date.Location()

	dayStart := time.Date(y, m, d, 0, 0, 0, 0, loc)

	return s.repo.GetSlotsByDate(mentorID, dayStart)
}

// for delete slots
func (s *service) DeleteSlot(userID uint, slotID uint) error {
	profileID, err := s.repo.GetProfileIDByUserID(userID)
	if err != nil {
		return errors.New("mentor profile not found")
	}
	slots, err := s.repo.GetSlotsByMentor(profileID)
	if err != nil {
		return err
	}

	var found *MentorSlot
	for i := range slots {
		if slots[i].ID == slotID {
			found = &slots[i]
			break
		}
	}
	if found == nil {
		return errors.New("slot not found")
	}
	return s.repo.DeleteSlot(slotID, profileID)
}
