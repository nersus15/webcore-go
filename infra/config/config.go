package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/webcore-go/webcore/app/helper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Memory   MemoryConfig   `mapstructure:"memory"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	PubSub   PubSubConfig   `mapstructure:"pubsub"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Others   map[string]ConfigObject
}

type AppConfig struct {
	Name              string          `mapstructure:"name"`
	Version           string          `mapstructure:"version"`
	Environment       string          `mapstructure:"environment"`
	Features          FeaturesConfig  `mapstructure:"features"`
	Logging           LoggingConfig   `mapstructure:"logging"`
	CORS              CORSConfig      `mapstructure:"cors"`
	RateLimit         RateLimitConfig `mapstructure:"rate_limit"`
	SecurityHeaders   bool            `mapstructure:"security_headers"`
	AdditionalHeaders []string        `mapstructure:"additional_headers"`
	Module            ModuleConfig    `mapstructure:"module"`
}

type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Max      int           `mapstructure:"max"`
	Window   time.Duration `mapstructure:"window"`
	Backend  string        `mapstructure:"backend"` // memory | redis
	FailOpen bool          `mapstructure:"fail_open"`
}

type FeaturesConfig struct {
	Logging   bool `mapstructure:"logging"`
	Recovery  bool `mapstructure:"recovery"`
	Metrics   bool `mapstructure:"metrics"`
	Tracing   bool `mapstructure:"tracing"`
	Profiling bool `mapstructure:"profiling"`
}

type LoggingConfig struct {
	Level  string              `mapstructure:"level"`
	Format string              `mapstructure:"format"`
	Output string              `mapstructure:"output"`
	Remote RemoteLoggingConfig `mapstructure:"remote"`
}

type RemoteLoggingConfig struct {
	Driver          string                    `mapstructure:"driver"`
	Uri             string                    `mapstructure:"uri"`
	SampleRate      int                       `mapstructure:"sample_rate"`
	Structured      bool                      `mapstructure:"structured"`
	DefaultTags     map[string]string         `mapstructure:"default_tags"`
	DefaultContexts map[string]map[string]any `mapstructure:"default_contexts"`
}

type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	PathPrefix   string        `mapstructure:"path"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ProxyHeader  string        `mapstructure:"proxy_header"` // header berisi IP klien dari proxy, kosong = alamat koneksi
}

type DatabaseConfig struct {
	Driver          string            `mapstructure:"driver"` // supported: "postgres", "mysql", "sqlite", "mongodb"
	Uri             string            `mapstructure:"uri"`
	Scheme          string            `mapstructure:"scheme"`
	Host            string            `mapstructure:"host"`
	Port            int               `mapstructure:"port"`
	User            string            `mapstructure:"user"`
	Password        string            `mapstructure:"password"`
	Name            string            `mapstructure:"name"`
	SchemaName      string            `mapstructure:"schema_name"`
	SSLMode         string            `mapstructure:"ssl_mode"`
	Attributes      map[string]string `mapstructure:"attributes"` // Additional connection parameters
	MaxIdleConns    int               `mapstructure:"max_idle_conns"`
	MaxOpenConns    int               `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration     `mapstructure:"conn_max_lifetime"`
	SlaveHosts      []DatabaseConfig  `mapstructure:"slave_hosts"`
}

type MemoryConfig struct {
	Enabled   bool          `mapstructure:"enabled"`
	Limit     int           `mapstructure:"limit"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

type RedisConfig struct {
	Host       string        `mapstructure:"host"`
	Port       int           `mapstructure:"port"`
	Password   string        `mapstructure:"password"`
	DB         int           `mapstructure:"db"`
	SlaveHosts []RedisConfig `mapstructure:"slave_hosts"`
}

type KafkaConfig struct {
	Enabled                     bool     `mapstructure:"enabled"`
	Brokers                     []string `mapstructure:"brokers"`
	GroupID                     string   `mapstructure:"group_id"`
	Topics                      []string `mapstructure:"topics"`
	AutoOffsetReset             string   `mapstructure:"offset_reset"`
	AutoCommit                  bool     `mapstructure:"auto_commit"`
	AutoOffsetStore             bool     `mapstructure:"auto_offset_store"`
	PartitionAssignmentStrategy string   `mapstructure:"assignment_strategy"`
}

type GoogleCredential struct {
	Type                    string `mapstructure:"type" json:"type"`
	ProjectID               string `mapstructure:"project_id" json:"project_id"`
	PrivateKeyID            string `mapstructure:"private_key_id" json:"private_key_id"`
	PrivateKey              string `mapstructure:"private_key" json:"private_key"`
	ClientEmail             string `mapstructure:"client_email" json:"client_email"`
	ClientID                string `mapstructure:"client_id" json:"client_id"`
	AuthURI                 string `mapstructure:"auth_uri" json:"auth_uri"`
	TokenURI                string `mapstructure:"token_uri" json:"token_uri"`
	AuthProviderX509CertURL string `mapstructure:"auth_provider_x509_cert_url" json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `mapstructure:"client_x509_cert_url" json:"client_x509_cert_url"`
	UniverseDomain          string `mapstructure:"universe_domain" json:"universe_domain"`
}

type PubSubConfig struct {
	Driver          string            `mapstructure:"driver"` // gpubsub, rabbitmq, awspubsub
	ProjectID       string            `mapstructure:"project_id"`
	Subscription    string            `mapstructure:"subscription"`
	Topic           string            `mapstructure:"topic"`
	NoCredentials   bool              `mapstructure:"no_credentials"`
	CredentialsPath string            `mapstructure:"credentials"`
	Credentials     *GoogleCredential `mapstructure:"credentials_data"`
	Consumer        ConsumerConfig    `mapstructure:"consumer"`
	Producer        ProducerConfig    `mapstructure:"producer"`
}

type ConsumerConfig struct {
	MaxMessagesPerPull    int               `mapstructure:"maxmessages"`
	SleepTimeBetweenPulls time.Duration     `mapstructure:"sleeptime"`
	AcknowledgeTimeout    time.Duration     `mapstructure:"acktimeout"`
	RetryCount            int               `mapstructure:"retrycount"`
	RetryDelay            time.Duration     `mapstructure:"retrydelay"`
	FlowControl           FlowControlConfig `mapstructure:"flowcontrol"`
}

type FlowControlConfig struct {
	Enabled                bool  `mapstructure:"enabled"`
	MaxOutstandingMessages int   `mapstructure:"maxmessages"`
	MaxOutstandingBytes    int64 `mapstructure:"maxbytes"`
}

type ProducerConfig struct {
	EnableMessageOrdering bool              `mapstructure:"enableordering"`
	BatchSize             int               `mapstructure:"batchsize"`
	MessageAttributes     map[string]string `mapstructure:"attributes"`
}

type AuthConfig struct {
	Directory    string            `mapstructure:"directory"` // e.g., "RBAC", "ABAC"
	Control      string            `mapstructure:"control"`   // e.g., "RBAC", "ABAC"
	Store        string            `mapstructure:"store"`     // e.g., "yaml", "db"
	Type         string            `mapstructure:"type"`      // e.g., "jwt", "apikey"
	Session      AuthSessionConfig `mapstructure:"session"`
	SecretKey    string            `mapstructure:"secret_key"`
	APIKeyHeader string            `mapstructure:"api_key_header"` // Header name for API key (default: "X-API-Key")
	APIKeyPrefix string            `mapstructure:"api_key_prefix"` // Optional prefix for API key validation
}

type AuthSessionConfig struct {
	Backend      string        `mapstructure:"backend"`      // e.g., "cache", "file", "db"
	ExpiresIn    time.Duration `mapstructure:"expires_in"`   // In seconds
	RefreshIn    time.Duration `mapstructure:"refresh_in"`   // In seconds
	ContentType  string        `mapstructure:"content_type"` // e.g., "application/json", "application/x-www-form-urlencoded"
	UsernameKey  string        `mapstructure:"username_key"`
	PasswordKey  string        `mapstructure:"password_key"`
	PublicClient bool          `mapstructure:"public_client"` // Set True if the client is for Browser, so the refresh token will read from coockie not request body
	CookieName   string        `mapstructure:"cookie_name"`   // cookie name for refresh_token
}

type ModuleConfig struct {
	Disabled []string `mapstructure:"disabled"`
	BasePath string   `mapstructure:"base_path"`
}

func (c *Config) GetOthers() map[string]ConfigObject {
	return c.Others
}

func (c *Config) AddOtherItem(key string, item ConfigObject) {
	ConfigAddItem(&c.Others, key, item)
}

func (c *Config) GetOtherItem(key string) (ConfigObject, bool) {
	if c.Others == nil {
		return nil, false
	}

	item, ok := c.Others[key]
	return item, ok
}

func (c *Config) GetFiberConfig(errorHandler fiber.ErrorHandler) fiber.Config {
	cfg := fiber.Config{
		ReadTimeout:   c.Server.ReadTimeout,
		WriteTimeout:  c.Server.WriteTimeout,
		CaseSensitive: true,
		StrictRouting: true,
		ErrorHandler:  errorHandler,

		// For faster encoder/decoder use go-json that wrap in helper.JSONMarshal and helper.JSONUnmarshal
		JSONEncoder: helper.JSONMarshal,
		JSONDecoder: helper.JSONUnmarshal,
	}

	if c.Server.ProxyHeader != "" {
		cfg.ProxyHeader = c.Server.ProxyHeader
		cfg.EnableIPValidation = true
	}

	return cfg
}
