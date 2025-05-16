package model

// Request body untuk membuat question link
type QuestionLinkRequest struct {
	QuestionLinkQuestionID int `json:"question_link_question_id"`
	QuestionLinkTargetType int `json:"question_link_target_type"` // 1=quiz, 2=evaluation, ...
	QuestionLinkTargetID   int `json:"question_link_target_id"`
}
