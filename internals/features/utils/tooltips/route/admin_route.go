package route

import (
	"quizku/internals/constants"
	tooltipController "quizku/internals/features/utils/tooltips/controller"
	authMiddleware "quizku/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TooltipAdminRoutes(router fiber.Router, db *gorm.DB) {
	tooltipCtrl := tooltipController.NewTooltipsController(db)
	tooltipInjectCtrl := tooltipController.NewTooltipInjectController(db)

	// 🔐 Route admin/teacher/owner
	adminRoutes := router.Group("/tooltip",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorNonUser("tooltip"),
			constants.TeacherAndAbove,
		),
	)

	adminRoutes.Get("/", tooltipCtrl.GetAllTooltips)
	adminRoutes.Post("/get-tooltips-id", tooltipCtrl.GetTooltipsID)
	adminRoutes.Post("/create-tooltips", tooltipCtrl.CreateTooltip)
	adminRoutes.Put("/:id", tooltipCtrl.UpdateTooltip)
	adminRoutes.Delete("/:id", tooltipCtrl.DeleteTooltip)

	// 💉 Route khusus untuk inject tooltip ID dalam teks
	adminRoutes.Post("/inject", tooltipInjectCtrl.InjectTooltipIDs)
}