package mentordiscovery




type Service struct{
	repo *Repository
}


func NewService(r *Repository)*Service{
	return &Service{repo: r}
}

func(s Service)GetMentors(filter *MentorFilter)([]MentorCard,error){
	
	//paginatin default
	if filter.Page<=0{
		filter.Page=1
	}
	if filter.Limit<=0{
		filter.Limit=10
	}
	filter.Offset = (filter.Page-1)*filter.Limit

	return s.repo.FindAllMentors(*filter)
}
