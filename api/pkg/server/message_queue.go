package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const Queue_name = "direct_agents"
const Questions_channel = "agents_questions_queue"
const Responses_channel = "agents_responses_queue"

// This is responsible for receiving commands from the server
func (scm *Connection) MessageQueueListen(listenChannel *amqp.Channel, queueName string) {
	messages, err := listenChannel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listenChannel.Close()

	for msg := range messages {
		var queue_response MQResponse

		err := json.Unmarshal(msg.Body, &queue_response)
		if err != nil {
			continue
		}

		scm.ParseAnswer(queue_response)
	}
}

func (scm *Connection) SendQuestion(sendChannel *amqp.Channel, sendQueueName string, question []byte) {
	err := sendChannel.Publish(
		Queue_name,
		sendQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        question,
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

	scm.MessageQueueListen(scm.AnswersChannel, responseQueueName)

	return nil
}

func (scm *Connection) DisconnectMQ() {
	if scm.MQConnection != nil && !scm.MQConnection.IsClosed() {
		scm.MQConnection.Close()
	}
}
