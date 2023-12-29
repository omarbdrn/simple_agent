package server

type MQQuestion struct {
	QuestionID string `json:"question_id"`
	Question   string `json:"question"`
}

type MQResponse struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}
