package mentordiscovery

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB)*Repository{
	return &Repository{db: db}
}

//repo function for find all mentors fitterr,search,sorting with pagination
func (r *Repository) FindAllMentors(f MentorFilter) ([]MentorCard, error) {
    var mentors []MentorCard

    q := r.db.
        Table("mentor_profiles").
        Select(`
            mentor_profiles.id as mentor_profile_id,
            users.id as user_id,
            users.name,
            users.profile_pic_url,
            mentor_profiles.bio,
            mentor_profiles.skills,
            mentor_profiles.hourly_price,
            mentor_profiles.rating_avg,
            mentor_profiles.experience_years
        `).
        Joins("join users on users.id = mentor_profiles.user_id").
        Where("users.role = ?", "mentor")

    if f.Skill != "" {
        q = q.Where("mentor_profiles.skills ILIKE ?", "%"+f.Skill+"%")
    }

    if f.Search != "" {
        q = q.Where("users.name ILIKE ?", "%"+f.Search+"%")
    }

    if f.MinPrice > 0 {
        q = q.Where("mentor_profiles.hourly_price >= ?", f.MinPrice)
    }

    if f.MaxPrice > 0 {
        q = q.Where("mentor_profiles.hourly_price <= ?", f.MaxPrice)
    }

    if f.MinRating > 0 {
        q = q.Where("mentor_profiles.rating_avg >= ?", f.MinRating)
    }

    // Sorting
    orderMap := map[string]string{
        "price_low":      "mentor_profiles.hourly_price ASC",
        "price_high":     "mentor_profiles.hourly_price DESC",
        "rating_high":    "mentor_profiles.rating_avg DESC",
        "experience_high":"mentor_profiles.experience_years DESC",
    }

    if order, ok := orderMap[f.Sort]; ok {
        q = q.Order(order)
    } else {
        q = q.Order("mentor_profiles.id DESC")
    }

    err := q.
        Limit(f.Limit).
        Offset(f.Offset).
        Scan(&mentors).Error

    return mentors, err
}

