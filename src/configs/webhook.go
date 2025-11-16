package configs

import (
	"context"
	"dainxor/atv/logger"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type webhookNS struct {
	configs map[string]string
}

var WebHooks webhookNS

func init() {
	urlsString, ok := os.LookupEnv("WH_URLS")

	if !ok {
		logger.Error("Failed to initialize the webhook")
		return
	}

	re := regexp.MustCompile(`\{([^,]+),([^}]+)\}`)
	matches := re.FindAllStringSubmatch(urlsString, -1)
	WebHooks.configs = make(map[string]string)

	for _, m := range matches {
		WebHooks.configs[strings.TrimSpace(m[1])] = strings.TrimSpace(m[2])
	}

	fmt.Printf("%+v\n", WebHooks)

	WebHooks.setupAMQP("cloudamqp")
}
func (webhookNS) Close() {}

func (w *webhookNS) setupAMQP(name string) error {
	url, ok := w.configs[name]

	if !ok {
		logger.Error("There is no url under the name", name)
		return errors.New("Name does not exists")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ:", err)
		return errors.New("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel:", err)
		return errors.New("Failed to open a channel: " + err.Error())
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"fatv-exchange", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		logger.Error("Failed to declare an exchange:", err)
		return errors.New("Failed to declare an exchange: " + err.Error())
	}

	q, err := ch.QueueDeclare(
		"fatv-app",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to declare the queue:", err)
		return errors.New("Failed to declare the queue: " + err.Error())
	}

	err = ch.QueueBind(
		q.Name,
		"fatv.*", // routing key pattern
		"fatv-exchange",
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to bind the queue:", err)
		return errors.New("Failed to bind the queue: " + err.Error())
	}

	return nil
}

func TestWH(data any) {
	WebHooks.SendTo("cloudamqp", Payload{
		Event:  "Test",
		Data:   data,
		SentAt: time.Now(),
	})

}

type Payload struct {
	Event  string    `json:"event"`
	Data   any       `json:"data"`
	SentAt time.Time `json:"sent_at"`
}

func (w webhookNS) SendTo(name string, payload Payload) error {
	url := w.configs[name]

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ:", err)
		return errors.New("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel:", err)
		return errors.New("Failed to open a channel: " + err.Error())
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"fatv-exchange", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		logger.Error("Failed to declare an exchange:", err)
		return errors.New("Failed to declare an exchange: " + err.Error())
	}

	q, err := ch.QueueDeclare(
		"fatv-app",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to declare the queue:", err)
		return errors.New("Failed to declare the queue: " + err.Error())
	}

	err = ch.QueueBind(
		q.Name,
		"fatv.*", // routing key pattern
		"fatv-exchange",
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to bind the queue:", err)
		return errors.New("Failed to bind the queue: " + err.Error())
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal the payload:", err)
		return errors.New("Failed to marshal the payload: " + err.Error())
	}

	err = ch.PublishWithContext(
		context.Background(),
		"fatv-exchange", // exchange
		"fatv.test",     // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payloadBytes,
		},
	)
	if err != nil {
		logger.Error("Failed to publish a message:", err)
		return errors.New("Failed to publish a message: " + err.Error())
	}

	logger.Debugf(" [x] Sent %s", payloadBytes)

	return err
}
