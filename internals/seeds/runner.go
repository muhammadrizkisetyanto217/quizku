package seeds

import (
	categories "quizku/internals/seeds/lessons/categories"
	difficulties "quizku/internals/seeds/lessons/difficulties"
	subcategories "quizku/internals/seeds/lessons/subcategories"
	themes_or_levels "quizku/internals/seeds/lessons/themes_or_levels"
	units "quizku/internals/seeds/lessons/units"

	evaluations "quizku/internals/seeds/quizzes/evaluations"
	exams "quizku/internals/seeds/quizzes/exams"
	questions "quizku/internals/seeds/quizzes/questions"
	quizzes "quizku/internals/seeds/quizzes/quizzes"
	reading "quizku/internals/seeds/quizzes/readings"
	section_quizzes "quizku/internals/seeds/quizzes/section_quizzes"

	level "quizku/internals/seeds/progress/levels"
	rank "quizku/internals/seeds/progress/ranks"

	"gorm.io/gorm"
)

func RunAllSeeds(db *gorm.DB) {

	//* Category
	difficulties.SeedDifficultiesFromJSON(db, "internals/seeds/category/difficulty/data_difficulty.json")
	categories.SeedCategoriesFromJSON(db, "internals/seeds/category/category/data_category.json")
	subcategories.SeedSubcategoriesFromJSON(db, "internals/seeds/category/subcategory/data_subcategory.json")
	themes_or_levels.SeedThemesOrLevelsFromJSON(db, "internals/seeds/category/themes_or_levels/data_themes_or_levels.json")
	units.SeedUnitsFromJSON(db, "internals/seeds/category/units/data_units.json")

	//* User

	//* Quizzes
	evaluations.SeedEvaluationsFromJSON(db, "internals/seeds/quizzes/evaluations/data_evaluations.json")
	exams.SeedExamsFromJSON(db, "internals/seeds/quizzes/exams/data_exams.json")
	questions.SeedQuestionsFromJSON(db, "internals/seeds/quizzes/questions/data_questions.json")
	quizzes.SeedQuizzesFromJSON(db, "internals/seeds/quizzes/quizzes/data_quizzes.json")
	reading.SeedReadingsFromJSON(db, "internals/seeds/quizzes/readings/data_readings.json")
	section_quizzes.SeedSectionQuizzesFromJSON(db, "internals/seeds/quizzes/section_quizzes/data_section_quizzes.json")

	//* Progress
	level.SeedLevelRequirementsFromJSON(db, "internals/seeds/progress/levels/data_levels_requirements.json")
	rank.SeedRanksRequirementsFromJSON(db, "internals/seeds/progress/ranks/data_ranks_requirements.json")

}
