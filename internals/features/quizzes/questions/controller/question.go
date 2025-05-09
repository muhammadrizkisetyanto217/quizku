package controller

import (
	"fmt"
	"log"
	"quizku/internals/features/quizzes/questions/model"
	questionModel "quizku/internals/features/quizzes/questions/model"
	tooltipModel "quizku/internals/features/utils/tooltips/model"
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuizzesQuestionController struct {
	DB *gorm.DB
}

func NewQuestionController(db *gorm.DB) *QuizzesQuestionController {
	return &QuizzesQuestionController{DB: db}
}

// GET all quiz questions
func (qqc *QuizzesQuestionController) GetQuestions(c *fiber.Ctx) error {
	log.Println("[INFO] Fetching all quiz questions")

	var questions []questionModel.QuestionModel
	if err := qqc.DB.Find(&questions).Error; err != nil {
		log.Println("[ERROR] Failed to fetch quiz questions:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch quiz questions",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d quiz questions\n", len(questions))
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "All quiz questions fetched successfully",
		"total":   len(questions),
		"data":    questions,
	})
}

// GET single quiz question
func (qqc *QuizzesQuestionController) GetQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching quiz question by ID: %s\n", id)

	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, id).Error; err != nil {
		log.Println("[ERROR] Quiz question not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Quiz question not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Quiz question fetched successfully by ID",
		"data":    question,
	})
}

// GET quiz questions by quiz ID
func (qqc *QuizzesQuestionController) GetQuestionsByQuizID(c *fiber.Ctx) error {
	quizID := c.Params("quizId")
	log.Printf("[INFO] Fetching quiz questions linked to quiz ID: %s\n", quizID)

	var links []model.QuestionLink
	if err := qqc.DB.
		Where("target_type = ? AND target_id = ?", model.TargetTypeQuiz, quizID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for quiz_id %s: %v\n", quizID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links",
		})
	}

	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	if len(questionIDs) == 0 {
		log.Printf("[INFO] No questions linked to quiz_id %s\n", quizID)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "No questions found for this quiz",
			"total":   0,
			"data":    []any{},
		})
	}

	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch questions by IDs: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch questions",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d questions linked to quiz_id %s\n", len(questions), quizID)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Quiz questions fetched successfully",
		"total":   len(questions),
		"data":    questions,
	})
}

// GET quiz questions by evaluation ID
func (qqc *QuizzesQuestionController) GetQuestionsByEvaluationID(c *fiber.Ctx) error {
	evaluationID := c.Params("evaluationId")
	log.Printf("[INFO] Fetching evaluation questions linked to evaluation ID: %s\n", evaluationID)

	// Ambil data dari question_links
	var links []model.QuestionLink
	if err := qqc.DB.
		Where("target_type = ? AND target_id = ?", model.TargetTypeEvaluation, evaluationID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for evaluation_id %s: %v\n", evaluationID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links for evaluation",
		})
	}

	// Ambil question_id dari links
	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	if len(questionIDs) == 0 {
		log.Printf("[INFO] No questions linked to evaluation_id %s\n", evaluationID)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "No questions found for this evaluation",
			"total":   0,
			"data":    []any{},
		})
	}

	// Ambil question dari tabel questions
	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch questions by IDs: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch questions for evaluation",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d questions linked to evaluation_id %s\n", len(questions), evaluationID)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Evaluation questions fetched successfully",
		"total":   len(questions),
		"data":    questions,
	})
}

// Get quiz questions by exam ID
func (qqc *QuizzesQuestionController) GetQuestionsByExamID(c *fiber.Ctx) error {
	examID := c.Params("examId")
	log.Printf("[INFO] Fetching exam questions linked to exam ID: %s\n", examID)

	// Ambil data dari question_links dengan target_type = 3 (exam)
	var links []model.QuestionLink
	if err := qqc.DB.
		Where("target_type = ? AND target_id = ?", model.TargetTypeExam, examID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for exam_id %s: %v\n", examID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links for exam",
		})
	}

	// Ambil question_id dari links
	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	if len(questionIDs) == 0 {
		log.Printf("[INFO] No questions linked to exam_id %s\n", examID)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "No questions found for this exam",
			"total":   0,
			"data":    []any{},
		})
	}

	// Ambil question dari tabel questions
	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch questions by IDs: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch questions for exam",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d questions linked to exam_id %s\n", len(questions), examID)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Exam questions fetched successfully",
		"total":   len(questions),
		"data":    questions,
	})
}

// Get quiz questions by test ID
func (qqc *QuizzesQuestionController) GetQuestionsByTestID(c *fiber.Ctx) error {
	testID := c.Params("testId")
	log.Printf("[INFO] Fetching test_exam questions linked to test ID: %s\n", testID)

	// Ambil data dari question_links dengan target_type = 4 (test_exam)
	var links []model.QuestionLink
	if err := qqc.DB.
		Where("target_type = ? AND target_id = ?", model.TargetTypeTest, testID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for test_id %s: %v\n", testID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links for test_exam",
		})
	}

	// Ambil question_id dari links
	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	if len(questionIDs) == 0 {
		log.Printf("[INFO] No questions linked to test_id %s\n", testID)
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "No questions found for this test_exam",
			"total":   0,
			"data":    []any{},
		})
	}

	// Ambil question dari tabel questions
	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch questions by IDs: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch questions for test_exam",
		})
	}

	log.Printf("[SUCCESS] Retrieved %d questions linked to test_id %s\n", len(questions), testID)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Test exam questions fetched successfully",
		"total":   len(questions),
		"data":    questions,
	})
}

// POST create quiz question
func (qqc *QuizzesQuestionController) CreateQuestion(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create question(s)")

	var (
		single   questionModel.QuestionModel
		multiple []questionModel.QuestionModel
	)

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		// ✅ Input berupa array
		if err := c.BodyParser(&multiple); err != nil {
			log.Printf("[ERROR] Failed to parse array of questions: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON array"})
		}

		if len(multiple) == 0 {
			log.Println("[ERROR] Empty question array")
			return c.Status(400).JSON(fiber.Map{"error": "Array of questions is empty"})
		}

		// Simpan batch
		if err := qqc.DB.Create(&multiple).Error; err != nil {
			log.Printf("[ERROR] Failed to insert questions: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create questions"})
		}

		log.Printf("[SUCCESS] Inserted %d questions", len(multiple))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple questions created successfully",
			"data":    multiple,
		})
	}

	// ✅ Input berupa objek tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single question: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if err := qqc.DB.Create(&single).Error; err != nil {
		log.Printf("[ERROR] Failed to create quiz question: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create question"})
	}

	log.Printf("[SUCCESS] Question created with ID: %d", single.ID)
	return c.Status(201).JSON(fiber.Map{
		"message": "Question created successfully",
		"data":    single,
	})
}

// PUT update quiz question
func (qqc *QuizzesQuestionController) UpdateQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating quiz question with ID: %s\n", id)

	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, id).Error; err != nil {
		log.Println("[ERROR] Quiz question not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Quiz question not found",
		})
	}

	if err := c.BodyParser(&question); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request",
		})
	}

	question.QuestionAnswer = pq.StringArray(question.QuestionAnswer)

	if err := qqc.DB.Save(&question).Error; err != nil {
		log.Println("[ERROR] Failed to update quiz question:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update quiz question",
		})
	}

	log.Printf("[SUCCESS] Quiz question with ID %s updated\n", id)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Quiz question updated successfully",
		"data":    question,
	})
}

// DELETE quiz question
func (qqc *QuizzesQuestionController) DeleteQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting quiz question with ID: %s\n", id)

	if err := qqc.DB.Delete(&questionModel.QuestionModel{}, id).Error; err != nil {
		log.Println("[ERROR] Failed to delete quiz question:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to delete quiz question",
		})
	}

	log.Printf("[SUCCESS] Quiz question with ID %s deleted\n", id)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": fmt.Sprintf("Quiz question with ID %s deleted successfully", id),
	})
}

func (qqc *QuizzesQuestionController) MarkKeywords(text string, tooltips []tooltipModel.Tooltip) string {
	log.Printf("[DEBUG] Original question text: %s\n", text)

	for _, tooltip := range tooltips {
		keyword := tooltip.Keyword
		keywordID := strconv.Itoa(int(tooltip.ID))

		// Gunakan regex case-insensitive, tapi tetap pertahankan casing asli match-nya
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return match + "=" + keywordID
		})

		log.Printf("[DEBUG] Replacing all '%s' with '%s' in text", keyword, keyword+"="+keywordID)
	}

	log.Printf("[DEBUG] Modified question text: %s\n", text)
	return text
}

func (qqc *QuizzesQuestionController) GetQuestionWithTooltipsMarked(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching marked quiz question with tooltips, ID: %s\n", id)

	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, id).Error; err != nil {
		log.Println("[ERROR] Quiz question not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Quiz question not found",
		})
	}

	var tooltips []tooltipModel.Tooltip
	if len(question.TooltipsID) > 0 {
		if err := qqc.DB.Where("id = ANY(?)", pq.Array(question.TooltipsID)).Find(&tooltips).Error; err != nil {
			log.Println("[ERROR] Failed to fetch tooltips:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to fetch tooltips",
			})
		}
	}

	// Tandai keyword di berbagai bagian teks
	markedText := qqc.MarkKeywords(question.QuestionText, tooltips)
	markedExplain := qqc.MarkKeywords(question.ExplainQuestion, tooltips)
	markedAnswer := qqc.MarkKeywords(question.AnswerText, tooltips)
	markedParagraph := qqc.MarkKeywords(question.ParagraphHelp, tooltips)

	log.Printf("[SUCCESS] Marked and fetched quiz question ID: %s\n", id)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Quiz question with marked tooltips fetched successfully",
		"quiz_question": fiber.Map{
			"id":               question.ID,
			"question_text":    markedText,
			"question_answer":  question.QuestionAnswer,
			"question_correct": question.QuestionCorrect,
			"tooltips_id":      question.TooltipsID,
			"status":           question.Status,
			"paragraph_help":   markedParagraph,
			"explain_question": markedExplain,
			"answer_text":      markedAnswer,
			"created_at":       question.CreatedAt,
			"updated_at":       question.UpdatedAt,
		},
		"tooltips": tooltips,
	})
}

// GET quiz question + tooltips
func (qqc *QuizzesQuestionController) GetQuestionWithTooltips(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching quiz question with tooltips, ID: %s\n", id)

	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, id).Error; err != nil {
		log.Println("[ERROR] Quiz question not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Quiz question not found",
		})
	}

	var tooltips []tooltipModel.Tooltip
	if len(question.TooltipsID) > 0 {
		if err := qqc.DB.Where("id = ANY(?)", pq.Array(question.TooltipsID)).Find(&tooltips).Error; err != nil {
			log.Println("[ERROR] Failed to fetch tooltips:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to fetch tooltips",
			})
		}
	}

	log.Printf("[SUCCESS] Retrieved quiz question and tooltips for ID: %s\n", id)
	return c.JSON(fiber.Map{
		"status":        true,
		"message":       "Quiz question and tooltips fetched successfully",
		"quiz_question": question,
		"tooltips":      tooltips,
	})
}

// GET only tooltips by quiz question ID
func (qqc *QuizzesQuestionController) GetOnlyQuestionTooltips(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching tooltips only for quiz question ID: %s\n", id)

	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, id).Error; err != nil {
		log.Println("[ERROR] Quiz question not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Quiz question not found",
		})
	}

	var tooltips []tooltipModel.Tooltip
	if len(question.TooltipsID) > 0 {
		if err := qqc.DB.Where("id = ANY(?)", pq.Array(question.TooltipsID)).Find(&tooltips).Error; err != nil {
			log.Println("[ERROR] Failed to fetch tooltips:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to fetch tooltips",
			})
		}
	}

	log.Printf("[SUCCESS] Retrieved tooltips only for quiz question ID: %s\n", id)
	return c.JSON(fiber.Map{
		"status":   true,
		"message":  "Tooltips fetched successfully",
		"tooltips": tooltips,
	})
}
