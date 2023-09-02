package lib

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/raman-vhd/arvan-challenge/internal/model"
)

type EventProducer struct {
	p *kafka.Producer
}

func NewEventProducer(env Env) EventProducer {
	conf := kafka.ConfigMap{
		"bootstrap.servers":  env.KafkaBootstrapServer,
		"acks":               "1",
		"request.timeout.ms": "5000",
	}
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v\n", err)
	}

	return EventProducer{
        p: p,
	}
}

func (p EventProducer) ProduceEvent(data model.Data, topic string) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.p.Produce(
		&kafka.Message{
			Value:          dataJSON,
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		},
		nil,
	)

	return err
}

// type EventConsumer struct {
// 	*kafka.Consumer
// }

// func NewEventConsumer(env Env) EventConsumer {
// 	conf := kafka.ConfigMap{
// 		"bootstrap.servers": env.KafkaBootstrapServer,
// 		"group.id":          "data-handler",
// 		// "auto.offset.reset": "smallest",
// 	}
// 	c, err := kafka.NewConsumer(&conf)
// 	if err != nil {
// 		log.Fatalf("failed to create kafka consumer: %v\n", err)
// 	}

// 	return EventConsumer{
// 		Consumer: c,
// 	}
// }

// func (c EventConsumer) Consume(topics []string) error {
// 	err = c.SubscribeTopics(topics, nil)
//     if err != nil {
//         return err
//     }

//     run := true
// 	for run == true {
// 		ev := c.Poll(500)
// 		switch e := ev.(type) {
// 		case *kafka.Message:
//             return 
// 		case kafka.Error:
//             return e
// 		}
// 	}

// 	c.Close()
// }
