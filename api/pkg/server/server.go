package server

import (
	"github.com/streadway/amqp"
)

type Connection struct {
	MQConnection     *amqp.Connection
	QuestionsChannel *amqp.Channel
	AnswersChannel   *amqp.Channel
}

func NewServer() *Connection {
	return &Connection{}
}
