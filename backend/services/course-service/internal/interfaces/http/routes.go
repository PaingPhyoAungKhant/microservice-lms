package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/paingphyoaungkhant/asto-microservice/services/course-service/cmd/docs"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(
	router *gin.Engine,
	categoryHandler *handlers.CategoryHandler,
	courseHandler *handlers.CourseHandler,
	courseOfferingHandler *handlers.CourseOfferingHandler,
	courseSectionHandler *handlers.CourseSectionHandler,
	sectionModuleHandler *handlers.SectionModuleHandler,
	logger *logger.Logger,
) {
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	// Swagger documentation
	router.GET("/api/v1/courses/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api/v1")
	{
		categoryRoutes := api.Group("/categories")
		{
			categoryRoutes.GET("", categoryHandler.FindCategory)
			categoryRoutes.GET("/:id", categoryHandler.GetCategory)
			categoryRoutes.POST("", categoryHandler.CreateCategory)
			categoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
			categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		courseRoutes := api.Group("/courses")
		{
			courseRoutes.GET("", courseHandler.FindCourse)
			courseRoutes.GET("/:id", courseHandler.GetCourse)
			courseRoutes.GET("/:id/details", courseHandler.GetCourseWithDetails)
			courseRoutes.POST("", courseHandler.CreateCourse)
			courseRoutes.PUT("/:id", courseHandler.UpdateCourse)
			courseRoutes.DELETE("/:id", courseHandler.DeleteCourse)

			// Course offerings
			courseRoutes.POST("/:course_id/offerings", courseOfferingHandler.CreateCourseOffering)
		}

		// Course offerings
		courseOfferingRoutes := api.Group("/course-offerings")
		{
			courseOfferingRoutes.GET("", courseOfferingHandler.FindCourseOffering)
			courseOfferingRoutes.GET("/:offering_id", courseOfferingHandler.GetCourseOffering)
			courseOfferingRoutes.PUT("/:offering_id", courseOfferingHandler.UpdateCourseOffering)
			courseOfferingRoutes.DELETE("/:offering_id", courseOfferingHandler.DeleteCourseOffering)
			courseOfferingRoutes.POST("/:offering_id/instructors", courseOfferingHandler.AssignInstructor)
			courseOfferingRoutes.DELETE("/:offering_id/instructors/:instructor_id", courseOfferingHandler.RemoveInstructor)
		}

		// Course sections
		courseSectionRoutes := api.Group("/course-offerings")
		{
			courseSectionRoutes.POST("/:offering_id/sections", courseSectionHandler.CreateCourseSection)
			courseSectionRoutes.GET("/:offering_id/sections", courseSectionHandler.FindCourseSection)
			courseSectionRoutes.PUT("/:offering_id/sections/reorder", courseSectionHandler.ReorderCourseSections)
		}

		courseSectionUpdateRoutes := api.Group("/course-sections")
		{
			courseSectionUpdateRoutes.GET("/:section_id", courseSectionHandler.GetCourseSection)
			courseSectionUpdateRoutes.PUT("/:section_id", courseSectionHandler.UpdateCourseSection)
			courseSectionUpdateRoutes.DELETE("/:section_id", courseSectionHandler.DeleteCourseSection)
		}

		// Section modules
		sectionModuleRoutes := api.Group("/course-sections")
		{
			sectionModuleRoutes.POST("/:section_id/modules", sectionModuleHandler.CreateSectionModule)
			sectionModuleRoutes.GET("/:section_id/modules", sectionModuleHandler.FindSectionModule)
			sectionModuleRoutes.PUT("/:section_id/modules/reorder", sectionModuleHandler.ReorderSectionModules)
		}

		sectionModuleUpdateRoutes := api.Group("/section-modules")
		{
			sectionModuleUpdateRoutes.GET("/:module_id", sectionModuleHandler.GetSectionModule)
			sectionModuleUpdateRoutes.PUT("/:module_id", sectionModuleHandler.UpdateSectionModule)
			sectionModuleUpdateRoutes.DELETE("/:module_id", sectionModuleHandler.DeleteSectionModule)
		}
	}
}

