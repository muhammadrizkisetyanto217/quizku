package route

import (
	tooltipController "quizku/internals/features/utils/tooltips/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TooltipPublicRoutes(router fiber.Router, db *gorm.DB) {
	tooltipCtrl := tooltipController.NewTooltipsController(db)

	publicRoutes := router.Group("/tooltip")
	publicRoutes.Get("/:id", tooltipCtrl.GetTooltipByID) // bebas tanpa login
}
