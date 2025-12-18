package integration

import (
	"database/sql"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/infrastructure/persistence/postgres"
)

func SetupCategoryRepository(db *sql.DB) repositories.CategoryRepository {
	return postgres.NewPostgresCategoryRepository(db)
}

func SetupCourseRepository(db *sql.DB) repositories.CourseRepository {
	return postgres.NewPostgresCourseRepository(db)
}

func SetupCourseCategoryRepository(db *sql.DB) repositories.CourseCategoryRepository {
	return postgres.NewPostgresCourseCategoryRepository(db)
}

func SetupCourseOfferingRepository(db *sql.DB) repositories.CourseOfferingRepository {
	return postgres.NewPostgresCourseOfferingRepository(db)
}

func SetupCourseOfferingInstructorRepository(db *sql.DB) repositories.CourseOfferingInstructorRepository {
	return postgres.NewPostgresCourseOfferingInstructorRepository(db)
}

func SetupCourseSectionRepository(db *sql.DB) repositories.CourseSectionRepository {
	return postgres.NewPostgresCourseSectionRepository(db)
}

func SetupSectionModuleRepository(db *sql.DB) repositories.SectionModuleRepository {
	return postgres.NewPostgresSectionModuleRepository(db)
}

