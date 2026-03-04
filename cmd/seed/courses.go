package seed

import (
	"LevelUp_Hub_Backend/internal/modules/courses"

	"gorm.io/gorm"
)


var defaultCourses = []courses.Course{
	{
		Title:       "Full-Stack Go Developer",
		Description: "Master backend engineering using Gin, GORM, and PostgreSQL. Learn to build scalable microservices.",
		ImageURL:    "https://images.unsplash.com/photo-1515879218367-8466d910aaa4",
		Category:    "Software Development",
		Level:       "Intermediate",
		IsActive:    true,
	},
	{
		Title:       "Advanced React & Next.js",
		Description: "Deep dive into Server Components, Server Actions, and high-performance frontend architecture.",
		ImageURL:    "https://images.unsplash.com/photo-1633356122544-f134324a6cee",
		Category:    "Web Development",
		Level:       "Advanced",
		IsActive:    true,
	},
	{
		Title:       "UI/UX Design Fundamentals",
		Description: "Learn Figma, typography, and color theory to create stunning user interfaces and seamless experiences.",
		ImageURL:    "https://images.unsplash.com/photo-1586717791821-3f44a563eb4c",
		Category:    "Design",
		Level:       "Beginner",
		IsActive:    true,
	},
	{
		Title:       "Data Science with Python",
		Description: "Analyze large datasets and build predictive models using NumPy, Pandas, and Scikit-Learn.",
		ImageURL:    "https://images.unsplash.com/photo-1551288049-bbbda536339a",
		Category:    "Data Science",
		Level:       "Intermediate",
		IsActive:    true,
	},
	{
		Title:       "Cybersecurity Essentials",
		Description: "Understand the fundamentals of network security, ethical hacking, and threat mitigation.",
		ImageURL:    "https://images.unsplash.com/photo-1550751827-4bd374c3f58b",
		Category:    "IT & Security",
		Level:       "Beginner",
		IsActive:    true,
	},
	{
		Title:       "Digital Marketing Strategy",
		Description: "Master SEO, SEM, and social media marketing to grow businesses in the digital age.",
		ImageURL:    "https://images.unsplash.com/photo-1460925895917-afdab827c52f",
		Category:    "Marketing",
		Level:       "Intermediate",
		IsActive:    true,
	},
	{
		Title:       "Cloud Computing with AWS",
		Description: "Deploy and manage scalable infrastructure using EC2, S3, and Lambda functions.",
		ImageURL:    "https://images.unsplash.com/photo-1451187580459-43490279c0fa",
		Category:    "Cloud Computing",
		Level:       "Advanced",
		IsActive:    true,
	},
	{
		Title:       "Public Speaking & Leadership",
		Description: "Build confidence and communication skills to lead teams and present effectively.",
		ImageURL:    "https://images.unsplash.com/photo-1475721027785-f74eccf877e2",
		Category:    "Soft Skills",
		Level:       "Beginner",
		IsActive:    true,
	},
	{
		Title:       "Mobile App Development with Flutter",
		Description: "Build cross-platform applications for iOS and Android from a single codebase.",
		ImageURL:    "https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c",
		Category:    "Software Development",
		Level:       "Intermediate",
		IsActive:    true,
	},
	{
		Title:       "Blockchain & Smart Contracts",
		Description: "Learn Solidity and Ethereum development to build decentralized applications (DApps).",
		ImageURL:    "https://images.unsplash.com/photo-1639762681485-074b7f938ba0",
		Category:    "Software Development",
		Level:       "Advanced",
		IsActive:    true,
	},
}

// SeedCourses populates the database with initial course data
func SeedCourses(db *gorm.DB) error {
	for _, course := range defaultCourses {
		// We use Title as the unique identifier to avoid duplicates
		if err := db.Where("title = ?", course.Title).FirstOrCreate(&courses.Course{}, course).Error; err != nil {
			return err
		}
	}
	return nil
}