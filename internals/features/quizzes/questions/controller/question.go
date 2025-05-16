package controller

import (
	"fmt"
	"log"
	"quizku/internals/features/quizzes/questions/model"
	questionModel "quizku/internals/features/quizzes/questions/model"

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

// 🔍 GET /api/quiz-questions
// Mengambil semua soal quiz dari database, tanpa filter.
// Umumnya dipakai untuk admin melihat seluruh bank soal yang tersedia.
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

// 🔍 GET /api/quiz-questions/:id
// Mengambil detail satu soal quiz berdasarkan ID-nya.
// Dipakai saat frontend ingin menampilkan detail soal untuk diedit.
func (qqc *QuizzesQuestionController) GetQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Fetching quiz question by ID: %s\n", id)

	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, "question_id = ?", id).Error; err != nil {
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

// 📎 GET /api/quiz-questions/quiz/:quizId
// Mengambil semua soal yang terhubung ke quiz tertentu berdasarkan quiz_id.
// Ini memanfaatkan tabel relasi `question_links` dengan `target_type = "quiz"`.
// Berguna untuk menyusun urutan soal per kuis yang aktif.
func (qqc *QuizzesQuestionController) GetQuestionsByQuizID(c *fiber.Ctx) error {
	quizID := c.Params("quizId")
	log.Printf("[INFO] Fetching quiz questions linked to quiz ID: %s\n", quizID)

	var links []model.QuestionLink
	if err := qqc.DB.
		Where("question_link_target_type = ? AND question_link_target_id = ?", model.TargetTypeQuiz, quizID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for quiz_id %s: %v\n", quizID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links",
		})
	}

	// Ekstrak semua question_link_question_id dari link
	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionLinkQuestionID)
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
		Where("question_id IN ?", questionIDs).
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

// 🔍 GET /api/quiz-questions/evaluation/:evaluationId
// Mengambil semua soal yang terhubung ke sebuah Evaluation.
// Menggunakan `question_links` dengan target_type = 2 (Evaluation).
func (qqc *QuizzesQuestionController) GetQuestionsByEvaluationID(c *fiber.Ctx) error {
	evaluationID := c.Params("evaluationId")
	log.Printf("[INFO] Fetching evaluation questions linked to evaluation ID: %s\n", evaluationID)

	var links []model.QuestionLink
	if err := qqc.DB.
		Where("question_link_target_type = ? AND question_link_target_id = ?", model.TargetTypeEvaluation, evaluationID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for evaluation_id %s: %v\n", evaluationID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links for evaluation",
		})
	}

	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionLinkQuestionID)
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

	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("question_id IN ?", questionIDs).
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

// 🔍 GET /api/quiz-questions/exam/:examId
// Mengambil soal berdasarkan `exam_id`, dari table `question_links` target_type = 3 (Exam).
func (qqc *QuizzesQuestionController) GetQuestionsByExamID(c *fiber.Ctx) error {
	examID := c.Params("examId")
	log.Printf("[INFO] Fetching exam questions linked to exam ID: %s\n", examID)

	var links []model.QuestionLink
	if err := qqc.DB.
		Where("question_link_target_type = ? AND question_link_target_id = ?", model.TargetTypeExam, examID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for exam_id %s: %v\n", examID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links for exam",
		})
	}

	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionLinkQuestionID)
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

	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("question_id IN ?", questionIDs).
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

// 🔍 GET /api/quiz-questions/test/:testId
// Mengambil soal berdasarkan `test_exam`, dengan target_type = 4.
func (qqc *QuizzesQuestionController) GetQuestionsByTestID(c *fiber.Ctx) error {
	testID := c.Params("testId")
	log.Printf("[INFO] Fetching test_exam questions linked to test ID: %s\n", testID)

	var links []model.QuestionLink
	if err := qqc.DB.
		Where("question_link_target_type = ? AND question_link_target_id = ?", model.TargetTypeTest, testID).
		Find(&links).Error; err != nil {
		log.Printf("[ERROR] Failed to fetch question links for test_id %s: %v\n", testID, err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to fetch question links for test_exam",
		})
	}

	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionLinkQuestionID)
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

	var questions []questionModel.QuestionModel
	if err := qqc.DB.
		Where("question_id IN ?", questionIDs).
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

// ✅ POST /api/quiz-questions
// Menambahkan satu atau banyak pertanyaan kuis ke database.
// Bisa menerima input dalam bentuk object tunggal atau array of objects.
func (qqc *QuizzesQuestionController) CreateQuestion(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create question(s)")

	type QuestionWithLink struct {
		questionModel.QuestionModel
		TargetType int `json:"target_type"`
		TargetID   int `json:"target_id"`
	}

	var (
		single   QuestionWithLink
		multiple []QuestionWithLink
	)

	raw := c.Body()
	if len(raw) > 0 && raw[0] == '[' {
		// Input array
		if err := c.BodyParser(&multiple); err != nil {
			log.Printf("[ERROR] Failed to parse array of questions: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON array"})
		}

		if len(multiple) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Array of questions is empty"})
		}

		var questions []questionModel.QuestionModel
		for _, q := range multiple {
			questions = append(questions, q.QuestionModel)
		}

		if err := qqc.DB.Create(&questions).Error; err != nil {
			log.Printf("[ERROR] Failed to insert questions: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create questions"})
		}

		// Insert question_links
		for i, q := range multiple {
			link := model.QuestionLink{
				QuestionLinkQuestionID: int(questions[i].QuestionID),
				QuestionLinkTargetType: q.TargetType,
				QuestionLinkTargetID:   q.TargetID,
			}

			if err := qqc.DB.Create(&link).Error; err != nil {
				log.Printf("[WARNING] Created question but failed to link question_id %d: %v", questions[i].QuestionID, err)
			}
		}

		log.Printf("[SUCCESS] Inserted %d questions and links", len(questions))
		return c.Status(201).JSON(fiber.Map{
			"message": "Multiple questions created and linked successfully",
			"data":    questions,
		})
	}

	// Input tunggal
	if err := c.BodyParser(&single); err != nil {
		log.Printf("[ERROR] Failed to parse single question: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if err := qqc.DB.Create(&single.QuestionModel).Error; err != nil {
		log.Printf("[ERROR] Failed to create quiz question: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create question"})
	}

	// Auto-create question_link
	link := model.QuestionLink{
		QuestionLinkQuestionID: int(single.QuestionID),
		QuestionLinkTargetType: single.TargetType,
		QuestionLinkTargetID:   single.TargetID,
	}

	if err := qqc.DB.Create(&link).Error; err != nil {
		log.Printf("[WARNING] Created question but failed to link question_id %d: %v", single.QuestionID, err)
	}

	log.Printf("[SUCCESS] Question created with ID: %d and linked", single.QuestionID)
	return c.Status(201).JSON(fiber.Map{
		"message": "Question created and linked successfully",
		"data":    single.QuestionModel,
	})
}

// 📝 PUT /api/quiz-questions/:id
// Mengupdate data pertanyaan kuis berdasarkan ID.
// Wajib menyertakan data lengkap di body (bukan partial update).
func (qqc *QuizzesQuestionController) UpdateQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Updating quiz question with ID: %s\n", id)

	// 🔍 Cari pertanyaan berdasarkan ID
	var question questionModel.QuestionModel
	if err := qqc.DB.First(&question, "question_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Quiz question not found:", err)
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Quiz question not found",
		})
	}

	// 🧾 Parsing body ke struct question
	if err := c.BodyParser(&question); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request",
		})
	}

	// ✅ Pastikan tipe array-nya sesuai
	question.QuestionAnswerChoices = pq.StringArray(question.QuestionAnswerChoices)

	// 💾 Simpan perubahan
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

// ❌ DELETE /api/quiz-questions/:id
// Menghapus pertanyaan kuis berdasarkan ID
func (qqc *QuizzesQuestionController) DeleteQuestion(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Printf("[INFO] Deleting quiz question with ID: %s\n", id)

	// Hapus question_links terlebih dahulu
	if err := qqc.DB.
		Where("question_link_question_id = ?", id).
		Delete(&model.QuestionLink{}).Error; err != nil {
		log.Println("[ERROR] Failed to delete question links:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to delete related question links",
		})
	}

	// Hapus pertanyaan
	if err := qqc.DB.Delete(&questionModel.QuestionModel{}, "question_id = ?", id).Error; err != nil {
		log.Println("[ERROR] Failed to delete quiz question:", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to delete quiz question",
		})
	}

	log.Printf("[SUCCESS] Quiz question with ID %s and its links deleted\n", id)
	return c.JSON(fiber.Map{
		"status":  true,
		"message": fmt.Sprintf("Quiz question with ID %s and its links deleted successfully", id),
	})
}
