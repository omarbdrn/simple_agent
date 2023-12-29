package server

type MQQuestion struct {
	QuestionID string `json:"question_id"`
	Question   string `json:"question"`
}

type MQResponse struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type Answer struct {
	Confirm bool   `json:"confirm"`
	IPRange string `json:"ip_range"`
}
