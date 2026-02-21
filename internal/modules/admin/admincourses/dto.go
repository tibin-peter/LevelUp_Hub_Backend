package admincourses


type CourseFilter struct {
	Search   string `query:"search"`
	Category string `query:"category"`
	Level    string `query:"level"`

	Page  int `query:"page"`
	Limit int `query:"limit"`
}