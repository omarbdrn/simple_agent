package server

import (
	"encoding/json"
	"time"

	"github.com/omarbdrn/simple-api/pkg/database"
	"gorm.io/gorm"
)

type Answer struct {
	Confirm bool   `json:"confirm"`
	IPRange string `json:"ip_range"`
}

func (scm *Connection) ParseAnswer(answer MQResponse) {
	var answer_model Answer
	var question database.Question

	err := json.Unmarshal([]byte(answer.Answer), &answer_model)
	if err != nil {
		return
	}

	if answer_model.Confirm {
		db := database.GetDB()

		result := db.First(&question, "ip_range = ?", answer_model.IPRange)
		if result.Error != nil {
			return
		}

		question.Answered = true
		db.Save(&question)
	}
}

func CheckQuestion(ip_range string) {
	time.Sleep(60 * time.Second)
	var question database.Question

	db := database.GetDB()
	result := db.First(&question, "ip_range = ?", ip_range)
	if result.Error == nil {
		if question.Answered {
			return
		}
	}

	FreeIPRange(ip_range, db)
	db.Delete(&question)
}

func FreeIPRange(ip_range string, db *gorm.DB) {
	var IPRange database.IPRange

	result := db.First(&IPRange, "ip_range = ?", ip_range)
	if result.Error != nil {
		return
	}

	IPRange.Taken = false
	db.Save(&IPRange)
}
