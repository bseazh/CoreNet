package queue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Producer interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}

type Consumer interface {
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
}

type SaramaProducer struct {
	producer sarama.SyncProducer
}

type SaramaConsumer struct {
	group sarama.ConsumerGroup
}

func NewSaramaProducer(brokers []string, cfg *sarama.Config) (*SaramaProducer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}
	prepared := prepareProducerConfig(cfg)
	prod, err := sarama.NewSyncProducer(brokers, prepared)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}
	return &SaramaProducer{producer: prod}, nil
}

func (p *SaramaProducer) Publish(ctx context.Context, topic string, payload []byte) error {
	if p == nil || p.producer == nil {
		return errors.New("kafka producer is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(topic) == "" {
		return errors.New("topic is required")
	}
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(payload)}
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("publish to topic %s: %w", topic, err)
	}
	return nil
}

func (p *SaramaProducer) Close() error {
	if p == nil || p.producer == nil {
		return nil
	}
	err := p.producer.Close()
	p.producer = nil
	return err
}

func NewSaramaConsumer(brokers []string, groupID string, cfg *sarama.Config) (*SaramaConsumer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}
	if strings.TrimSpace(groupID) == "" {
		return nil, errors.New("consumer group id is required")
	}

	prepared := prepareConsumerConfig(cfg)
	group, err := sarama.NewConsumerGroup(brokers, groupID, prepared)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer group: %w", err)
	}
	return &SaramaConsumer{group: group}, nil
}

func (c *SaramaConsumer) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	if c == nil || c.group == nil {
		return errors.New("kafka consumer is not initialized")
	}
	if handler == nil {
		return errors.New("handler is required")
	}
	if strings.TrimSpace(topic) == "" {
		return errors.New("topic is required")
	}

	cgHandler := &consumerGroupHandler{handler: handler}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := c.group.Consume(ctx, []string{topic}, cgHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("consume topic %s: %w", topic, err)
		}
		if handlerErr := cgHandler.Err(); handlerErr != nil {
			return handlerErr
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		cgHandler.Reset()
	}
}

func (c *SaramaConsumer) Close() error {
	if c == nil || c.group == nil {
		return nil
	}
	err := c.group.Close()
	c.group = nil
	return err
}

type consumerGroupHandler struct {
	handler func([]byte) error
	mu      sync.Mutex
	err     error
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(msg.Value); err != nil {
			h.setErr(err)
			return err
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerGroupHandler) setErr(err error) {
	if err == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.err == nil {
		h.err = err
	}
}

func (h *consumerGroupHandler) Err() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.err
}

func (h *consumerGroupHandler) Reset() {
	h.mu.Lock()
	h.err = nil
	h.mu.Unlock()
}

func prepareProducerConfig(cfg *sarama.Config) *sarama.Config {
	if cfg == nil {
		cfg = sarama.NewConfig()
	}
	if cfg.Version == (sarama.KafkaVersion{}) {
		cfg.Version = sarama.V2_5_0_0
	}
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	if cfg.Producer.Retry.Max == 0 {
		cfg.Producer.Retry.Max = 5
	}
	if cfg.Producer.Compression == sarama.CompressionNone {
		cfg.Producer.Compression = sarama.CompressionSnappy
	}
	return cfg
}

func prepareConsumerConfig(cfg *sarama.Config) *sarama.Config {
	if cfg == nil {
		cfg = sarama.NewConfig()
	}
	if cfg.Version == (sarama.KafkaVersion{}) {
		cfg.Version = sarama.V2_5_0_0
	}
	cfg.Consumer.Return.Errors = true
	if cfg.Consumer.Group.Rebalance.GroupStrategies == nil {
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	}
	if cfg.Consumer.Offsets.Initial == 0 {
		cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	if cfg.Consumer.Group.Session.Timeout == 0 {
		cfg.Consumer.Group.Session.Timeout = 10 * time.Second
	}
	if cfg.Consumer.Group.Heartbeat.Interval == 0 {
		cfg.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	}
	if cfg.Consumer.Offsets.AutoCommit.Interval == 0 {
		cfg.Consumer.Offsets.AutoCommit.Interval = time.Second
	}
	return cfg
}
