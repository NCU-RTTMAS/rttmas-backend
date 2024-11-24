package mqtt

// import (
// 	"context"
// 	"fmt"
// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"rttmas-backend/config"
// 	"rttmas-backend/pkg/utils/logger"
// )

// var MqConn *amqp.Connection
// var MqChannel *amqp.Channel

// func Init() (*amqp.Connection, *amqp.Channel) {

// 	uri := fmt.Sprintf("amqp://%s:%s@%s", config.GetConfigValue("AMQP_USERNAME"), config.GetConfigValue("AMQP_PASSWORD"), config.GetConfigValue("AMQP_BROKER_URI"))
// 	// logger.Info(uri)
// 	logger.Debug("[AMQP] Connecting to AMQP:", uri)

// 	conn, err := amqp.Dial(uri)
// 	if err != nil {
// 		logger.Fatal(err.Error())
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		logger.Fatal(err.Error())
// 	}

// 	err = ch.ExchangeDeclare(
// 		"rttma", // name
// 		"topic", // kind/type: "direct", "fanout", "topic", "headers"
// 		true,    // durable
// 		false,   // auto-deleted
// 		false,   // internal
// 		false,   // no-wait
// 		nil,     // arguments
// 	)

// 	if err != nil {
// 		logger.Fatal(err.Error())
// 	}

// 	_, err = ch.QueueDeclare(
// 		"analysis_module_queue", // name
// 		true,                    // durable
// 		false,                   // auto-deleted
// 		false,                   // exclusive
// 		false,                   // no-wait
// 		nil,                     // arguments
// 	)

// 	if err != nil {
// 		logger.Fatal(err.Error())
// 	}
// 	err = ch.QueueBind(
// 		"analysis_module_queue",               // queue name
// 		"traffic.plate_recognition_reports.#", // routing key
// 		"rttma",                               // exchange
// 		false,                                 // no-wait
// 		nil,                                   // arguments
// 	)
// 	err = ch.QueueBind(
// 		"analysis_module_queue",           // queue name
// 		"traffic.user_location_reports.#", // routing key
// 		"rttma",                           // exchange
// 		false,                             // no-wait
// 		nil,                               // arguments
// 	)
// 	err = ch.QueueBind(
// 		"analysis_module_queue",            // queue name
// 		"traffic.vehicle_true_locations.#", // routing key
// 		"rttma",                            // exchange
// 		false,                              // no-wait
// 		nil,                                // arguments
// 	)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	MqChannel = ch
// 	CreateAMQPEndpoints()

// 	return MqConn, MqChannel

// }

// func Close() {
// 	MqChannel.Close()
// 	MqConn.Close()
// }

// func PublishToTopic(topic string, msg string) error {
// 	cxt := context.Background()
// 	return MqChannel.PublishWithContext(
// 		cxt,     // context
// 		"rttma", // exchange
// 		topic,   // key / routing key
// 		false,   // mandatory
// 		false,   // immediate

// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte(msg),
// 		},
// 	)
// }

// func CreateConsumer(queueName string, handler func(amqp.Delivery)) error {
// 	msgs, err := MqChannel.Consume(
// 		queueName,       // queue
// 		"test_receiver", // consumer
// 		true,            // auto-ack
// 		false,           // exclusive
// 		false,           // no-local
// 		false,           // no-wait
// 		nil,             // args
// 	)

// 	if err != nil {
// 		logger.Fatal("Failed to register a consumer: ", err)
// 		return err
// 	}

// 	// Creating a goroutine to handle messages from the queue
// 	go func() {
// 		for msg := range msgs {
// 			handler(msg)
// 		}
// 	}()

// 	logger.Info("[MQTT] Consumer registered for queue: " + queueName)
// 	return nil
// }
