package route

import (
	"quizku/internals/constants"
	surveyController "quizku/internals/features/users/survey/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SurveyAdminRoutes(api fiber.Router, db *gorm.DB) {
	surveyQuestionCtrl := surveyController.NewSurveyQuestionController(db)

	surveyRoutes := api.Group("/survey-questions",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorOwner("mengelola pertanyaan survei"),
			constants.OwnerAndAbove,
		),
	)
	surveyRoutes.Post("/", surveyQuestionCtrl.Create)
	surveyRoutes.Put("/:id", surveyQuestionCtrl.Update)
	surveyRoutes.Delete("/:id", surveyQuestionCtrl.Delete)
}
