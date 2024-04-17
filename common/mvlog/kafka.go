package mvlog

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafkaWriter(kafkaURL, topic string) {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Println(err)
			}
		},
	}
}

var records []kafka.Message

func BeginBatch() {
	records = make([]kafka.Message, 0)
}

func WriteBatch(key string, data []byte) {
	if kafkaWriter == nil {
		return
	}
	records = append(records, kafka.Message{

		Key:   []byte(key),
		Value: data,
	})
}

func EndBatch(ctx context.Context) error {
	if len(records) == 0 {
		return nil
	}

	err := kafkaWriter.WriteMessages(ctx, records...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
