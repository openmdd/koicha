package kafka

import (
	"context"
	"fmt"
	"sort"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// Client wraps franz-go operations used by koicha.
type Client struct {
	client *kgo.Client
}

// NewClient creates a Kafka client from koicha Kafka config.
func NewClient(cfg Config) (*Client, error) {
	if len(cfg.BootstrapServers) == 0 {
		return nil, fmt.Errorf("kafka bootstrap servers are required")
	}
	if cfg.Auth.Protocol != SecurityProtocolPlaintext {
		return nil, fmt.Errorf("unsupported kafka security protocol: %s", cfg.Auth.Protocol)
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.BootstrapServers...),
	}
	if cfg.ClientID != "" {
		opts = append(opts, kgo.ClientID(cfg.ClientID))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}
	return &Client{client: client}, nil
}

// Ping checks whether at least one Kafka broker is reachable.
func (c *Client) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx); err != nil {
		return fmt.Errorf("ping kafka: %w", err)
	}
	return nil
}

// ListTopics returns topics visible to the configured Kafka client.
func (c *Client) ListTopics(ctx context.Context) ([]Topic, error) {
	req := kmsg.NewMetadataRequest()
	req.Topics = nil
	req.AllowAutoTopicCreation = false

	resp, err := req.RequestWith(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("list kafka topics: %w", err)
	}
	if err := kerr.ErrorForCode(resp.ErrorCode); err != nil {
		return nil, fmt.Errorf("list kafka topics: %w", err)
	}

	topics := make([]Topic, 0, len(resp.Topics))
	for _, item := range resp.Topics {
		if err := kerr.ErrorForCode(item.ErrorCode); err != nil {
			return nil, fmt.Errorf("list kafka topic metadata: %w", err)
		}
		if item.Topic == nil || *item.Topic == "" {
			continue
		}
		topics = append(topics, Topic{Name: *item.Topic})
	}

	sort.Slice(topics, func(i, j int) bool {
		return topics[i].Name < topics[j].Name
	})
	return topics, nil
}

// Close releases franz-go client resources.
func (c *Client) Close() {
	c.client.Close()
}
