package admindashboard


type Service interface {
	Dashboard(filter string,)(*AdminDashboardDTO,error)
}

type service struct {
	repo Repository
}

func NewService(r Repository)Service {
	return &service{repo: r}
}

func (s *service) Dashboard(filter string,)(*AdminDashboardDTO,error){

	start,end := getRange(filter)

	students,err := s.repo.CountStudents()
	if err != nil { return nil,err }

	mentors,err := s.repo.CountMentors()
	if err != nil { return nil,err }

	courses,err := s.repo.CountOfCourses()
	if err != nil { return nil,err }

	revenue,err := s.repo.TotalPlatformRevenue(start,end)
	if err != nil { return nil,err }

	chart,err := s.repo.RevenueChart(start,end,filter)
	if err != nil { return nil,err }

	return &AdminDashboardDTO{
		TotalStudents: students,
		TotalMentors: mentors,
		ActiveCourses: courses,
		TotalRevenue: revenue,
		RevenueChart: chart,
	},nil
}