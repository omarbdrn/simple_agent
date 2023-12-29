package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/omarbdrn/simple_agent/pkg/global"
	"github.com/streadway/amqp"
)

const Queue_name = "direct_agents"
const Questions_channel = "agents_questions_queue"
const Responses_channel = "agents_responses_queue"

func (scm *Connection) ParseQuestion(msg amqp.Delivery, mqQuestion MQQuestion) {
	ip_range := mqQuestion.Question
	current_share := global.GetCurrentShare()
	if global.IsInArray(ip_range, current_share.CIDRs) {
		response := Answer{
			Confirm: true,
			IPRange: ip_range,
		}

		jsonified_answer, err := json.Marshal(response)
		if err != nil {
			return
		}

		answer := MQResponse{QuestionID: mqQuestion.QuestionID, Answer: string(jsonified_answer)}

		jsonified_answer, err = json.Marshal(answer)
		if err != nil {
			return
		}

		scm.SendAnswer(scm.AnswersChannel, Responses_channel, jsonified_answer)
		msg.Ack(true)
	}
}

func (scm *Connection) MessageQueueListen(listenChannel *amqp.Channel, queueName string) {
	messages, err := listenChannel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listenChannel.Close()

	for msg := range messages {
		var queue_question MQQuestion

		err := json.Unmarshal(msg.Body, &queue_question)
		if err != nil {
			continue
		}

		scm.ParseQuestion(msg, queue_question)
	}
}

func (scm *Connection) SendAnswer(sendChannel *amqp.Channel, sendQueueName string, answer []byte) {
	err := sendChannel.Publish(
		Queue_name,
		sendQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        answer,
		},
	)

	if err != nil {
		log.Println("Error publishing response:", err)
	}
}

func (scm *Connection) ConnectMQ(username string, password string, vhost string) error {
	var err error

	connectionString := fmt.Sprintf("amqp://%s:%s@localhost:5672/%s", username, password, vhost)

	scm.MQConnection, err = amqp.DialConfig(connectionString, amqp.Config{
		Vhost: vhost,
	})
	if err != nil {
		return err
	}

	scm.QuestionsChannel, err = scm.MQConnection.Channel()
	if err != nil {
		scm.MQConnection.Close()
		return err
	}

	scm.AnswersChannel, err = scm.MQConnection.Channel()
	if err != nil {
		return err
	}

	// We're Declaring Exchanges and Queue to make sure it exists with the required settings.
	err = scm.QuestionsChannel.ExchangeDeclare(Queue_name, "direct", true, false, false, false, nil)
	if err != nil {
		scm.MQConnection.Close()
		return err
	}

	// Declaring the Commands Queue
	queueName := Questions_channel
	_, err = scm.QuestionsChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		scm.MQConnection.Close()
		return err
	}

	// Declaring the Responses Queue
	responseQueueName := Responses_channel
	_, err = scm.AnswersChannel.QueueDeclare(responseQueueName, true, false, false, false, nil)
	if err != nil {
		scm.MQConnection.Close()
		return err
	}

	// Binding to Responses Queue to send responses
	err = scm.AnswersChannel.QueueBind(responseQueueName, responseQueueName, Queue_name, false, nil)
	if err != nil {
		scm.MQConnection.Close()
		return err
	}

	// Binding to Commands Queue to listen for commands
	err = scm.QuestionsChannel.QueueBind(queueName, queueName, Queue_name, false, nil)
	if err != nil {
		scm.MQConnection.Close()
		return err
	}

	scm.MessageQueueListen(scm.QuestionsChannel, queueName)

	return nil
}

func (scm *Connection) DisconnectMQ() {
	if scm.MQConnection != nil && !scm.MQConnection.IsClosed() {
		scm.MQConnection.Close()
	}
}
