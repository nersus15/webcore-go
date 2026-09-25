package port

import (
	"context"
	"time"
)

type DbMap map[string]any

type DbExpression struct {
	Expr string // simple column name or complex expression
	Op   string // only used when expression is column name
	Args []any
}

type IDatabase interface {
	Connector

	// Common methods that both databases should implement
	Ping(ctx context.Context) error

	// Get the underlying connection for specific operations
	GetConnection() any
	GetDriver() string
	GetName() string

	Count(ctx context.Context, table string, filter []DbExpression) (int64, error)
	Find(ctx context.Context, results any, table string, column []string, filter []DbExpression, sort map[string]int, limit int64, skip int64) error
	FindOne(ctx context.Context, result any, table string, column []string, filter []DbExpression, sort map[string]int) error
	InsertOne(ctx context.Context, table string, data any) (any, error)
	Update(ctx context.Context, table string, filter []DbExpression, data any) (int64, error)
	UpdateOne(ctx context.Context, table string, filter []DbExpression, data any) (int64, error)
	Delete(ctx context.Context, table string, filter []DbExpression) (int64, error)
	DeleteOne(ctx context.Context, table string, filter []DbExpression) (int64, error)

	StartMigration(ctx context.Context, service string, command string, dir string, args []string) error
}

// Generic for Memory Caching (ex: Redis, MemCached)
type ICacheMemory interface {
	Connector

	Set(key string, value any, ttl time.Duration) error
	Get(key string, outvalue any) bool
	Delete(key string) error
}

type IPubSub interface {
	Connector

	Publish(ctx context.Context, message any, attributes map[string]string) (string, error)
	RegisterReceiver(receiver PubSubReceiver)
	StartReceiving(ctx context.Context)
}

type IPubSubMessage interface {
	GetID() string
	GetData() []byte
	GetPublishTime() time.Time
	GetAttributes() map[string]string
}

type PubSubReceiver interface {
	Consume(ctx context.Context, messages []IPubSubMessage) (map[string]bool, error)
}

type IKafka interface {
	Connector

	Publish(ctx context.Context, topic string, message any) error
	Consume(ctx context.Context, topic string) (<-chan any, error)
}

type KafkaConsumer interface {
	Consume(ctx context.Context, message []byte) (bool, error)
}
