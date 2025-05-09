package model

// ✅ Struct untuk request-nya
type QuestionLinkRequest struct {
	QuestionID int `json:"question_id"`
	TargetType int `json:"target_type"` // 1=quiz, 2=exam, ...
	TargetID   int `json:"target_id"`
}
