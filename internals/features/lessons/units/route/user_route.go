package route

import (
	unitController "quizku/internals/features/lessons/units/controller"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UnitUserRoutes(api fiber.Router, db *gorm.DB) {
	unitCtrl := unitController.NewUnitController(db)
	unitNewsCtrl := unitController.NewUnitNewsController(db)
	userUnitCtrl := unitController.NewUserUnitController(db)

	unitRoutes := api.Group("/units")
	unitRoutes.Get("/", unitCtrl.GetUnits)
	unitRoutes.Get("/:id", unitCtrl.GetUnit)
	unitRoutes.Get("/themes-or-levels/:themesOrLevelId", unitCtrl.GetUnitByThemesOrLevels)

	unitNewsRoutes := api.Group("/units-news")
	unitNewsRoutes.Get("/", unitNewsCtrl.GetAll)
	unitNewsRoutes.Get("/:unit_id", unitNewsCtrl.GetByUnitID)
	unitNewsRoutes.Get("/detail/:id", unitNewsCtrl.GetByID)

	userUnitRoutes := api.Group("/user-units")
	userUnitRoutes.Get("/:user_id", userUnitCtrl.GetByUserID)
	userUnitRoutes.Get("/themes-or-levels/:themes_or_levels_id", userUnitCtrl.GetUserUnitsByThemesOrLevelsAndUserID)
}
   