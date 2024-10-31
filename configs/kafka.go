package configs

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers      []string
	MaxRetries   int
	RetryBackoff time.Duration
}

type KafkaManager struct {
	config    KafkaConfig
	producers map[string]*kafka.Writer
	mu        sync.RWMutex
}

func NewKafkaManager(config KafkaConfig) *KafkaManager {
	return &KafkaManager{
		config:    config,
		producers: make(map[string]*kafka.Writer),
	}
}

// GetProducer retrieves or creates a producer for the given topic
func (m *KafkaManager) GetProducer(topic string) *kafka.Writer {
	m.mu.RLock()
	producer, exists := m.producers[topic]
	m.mu.RUnlock()
	if exists {
		return producer
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	// Double-check to avoid duplicate creation
	if producer, exists := m.producers[topic]; exists {
		return producer
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(m.config.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	m.producers[topic] = writer
	return writer
}

// Close all producers when shutting down the service
func (m *KafkaManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, writer := range m.producers {
		writer.Close()
	}
}

type KafkaService struct {
	manager *KafkaManager
}

func NewKafkaService(config KafkaConfig) *KafkaService {
	manager := NewKafkaManager(config)
	return &KafkaService{manager: manager}
}

// PublishMessage allows any service to publish a message to any topic
func (ks *KafkaService) PublishMessage(ctx context.Context, topic, message string) error {
	producer := ks.manager.GetProducer(topic)

	msg := kafka.Message{
		Value: []byte(message),
	}
	var err error
	for i := 0; i < ks.manager.config.MaxRetries; i++ {
		err = producer.WriteMessages(ctx, msg)
		if err == nil {
			log.Printf("Message sent successfully to topic %s", topic)
			return nil
		}
		log.Printf("Retrying send to topic %s, attempt %d/%d", topic, i+1, ks.manager.config.MaxRetries)
		time.Sleep(ks.manager.config.RetryBackoff)
	}
	return err
}

func (ks *KafkaService) Close() {
	ks.manager.CloseAll()
}
