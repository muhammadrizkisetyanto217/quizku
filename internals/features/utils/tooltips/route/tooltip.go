package route

import (
	"quizku/internals/constants"
	tooltipController "quizku/internals/features/utils/tooltips/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TooltipRoute(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api", authMiddleware.AuthMiddleware(db))

	tooltipRoutes := api.Group("/tooltip",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorNonUser("tooltip"),
			constants.OwnerAndAbove,
		),
	)

	tooltipCtrl := tooltipController.NewTooltipsController(db)

	tooltipRoutes.Get("/", tooltipCtrl.GetAllTooltips)
	tooltipRoutes.Post("/get-tooltips-id", tooltipCtrl.GetTooltipsID)
	tooltipRoutes.Post("/create-tooltips", tooltipCtrl.CreateTooltip)
	tooltipRoutes.Put("/:id", tooltipCtrl.UpdateTooltip)
	tooltipRoutes.Delete("/:id", tooltipCtrl.DeleteTooltip)

}
