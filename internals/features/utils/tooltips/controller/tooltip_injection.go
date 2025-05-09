package controller

import (
	"fmt"
	"quizku/internals/features/utils/tooltips/model"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TooltipInjectController struct {
	DB *gorm.DB
}

func NewTooltipInjectController(db *gorm.DB) *TooltipInjectController {
	return &TooltipInjectController{DB: db}
}

func (tc *TooltipInjectController) InjectTooltipIDs(c *fiber.Ctx) error {
	var req struct {
		Text string `json:"text"`
	}

	if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Text) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"error":  "Invalid request body or empty text",
		})
	}

	original := req.Text
	processed := tc.replaceWithTooltipIDs(original)

	return c.JSON(fiber.Map{
		"status":   true,
		"original": original,
		"result":   processed,
	})
}

// Fungsi utama mengganti keyword[] menjadi keyword[ID] jika ditemukan di DB
func (tc *TooltipInjectController) replaceWithTooltipIDs(text string) string {
	re := regexp.MustCompile(`\b(\w+)\[\]`)
	matches := re.FindAllStringSubmatch(text, -1)

	seen := map[string]bool{}

	for _, match := range matches {
		keyword := match[1]
		if seen[keyword] {
			continue
		}
		seen[keyword] = true

		var tooltip model.Tooltip
		if err := tc.DB.Where("keyword = ?", keyword).First(&tooltip).Error; err == nil {
			from := keyword + "[]"
			to := keyword + "[" + fmt.Sprintf("%d", tooltip.ID) + "]"
			text = strings.Replace(text, from, to, 1)
		}
	}

	return text
}
