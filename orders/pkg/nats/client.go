package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type Config struct {
	Hosts  []string
	NKey   string
	IsTest bool
}

type Client struct {
	Conn *nats.Conn
}

func NewClient(ctx context.Context, hosts []string, nkey string, isTest bool) (*Client, error) {
	options := []nats.Option{
		nats.Name("your-service-name"),
		nats.Timeout(10 * time.Second),
		nats.ReconnectWait(5 * time.Second),
		nats.MaxReconnects(10),
	}

	// Если это тестовое окружение, добавляем опцию NoEcho
	if isTest {
		options = append(options, nats.NoEcho())
	}

	// Подключение к NATS
	nc, err := nats.Connect(hosts[0], options...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &Client{
		Conn: nc,
	}, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
