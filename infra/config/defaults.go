package config

func (c *Config) SetDefaults() map[string]any {
	return map[string]any{
		// App
		"app.name":                            "webcore-go",
		"app.version":                         "1.0.0",
		"app.environment":                     "development",
		"app.features.logging":                true,
		"app.features.recovery":               false,
		"app.features.tracing":                false,
		"app.features.metrics":                false,
		"app.features.profiling":              false,
		"app.logging.level":                   "info",
		"app.logging.format":                  "json",
		"app.logging.output":                  "stdout",
		"app.logging.remote.driver":           "sentry",
		"app.logging.remote.uri":              "",
		"app.logging.remote.sample_rate":      1.0,
		"app.logging.remote.structured":       true,
		"app.logging.remote.default_tags":     map[string]string{},
		"app.logging.remote.default_contexts": map[string]map[string]any{},
		"app.security_headers":                false,
		"app.additional_headers":              []string{},
		"app.cors.allow_origins":              []string{"*"},
		"app.cors.allow_methods":              []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		"app.cors.allow_headers":              []string{"Origin", "Content-Type", "Accept", "Authorization"},
		"app.cors.expose_headers":             []string{"Content-Length"},
		"app.cors.allow_credentials":          true,
		"app.cors.max_age":                    "24h", // 24 hours
		"app.rate_limit.enabled":              false,
		"app.rate_limit.max":                  1000,
		"app.rate_limit.window":               "1m",
		"app.rate_limit.backend":              "memory",
		"app.rate_limit.fail_open":            true,
		"app.module.base_path":                "./libs",
		"app.module.disabled":                 []string{},

		// Server
		"server.host":          "0.0.0.0",
		"server.port":          7272,
		"server.path":          "/api",
		"server.read_timeout":  "30s",
		"server.write_timeout": "30s",
		"server.proxy_header":  "",

		// Auth
		"auth.directory":             ".",
		"auth.control":               "RBAC",
		"auth.store":                 "yaml",
		"auth.type":                  "jwt",
		"auth.secret_key":            "",
		"auth.session.backend":       "",
		"auth.session.expires_in":    "2h",  // 2 hours
		"auth.session.refresh_in":    "48h", // 2 days
		"auth.session.content_type":  "application/x-www-form-urlencoded",
		"auth.session.public_client": false,
		"auth.session.cookie_name":   "",
		"auth.session.username_key":  "",
		"auth.session.password_key":  "",
		"auth.api_key_header":        "X-API-Key",
		"auth.api_key_prefix":        "",

		// Database
		"database.driver":            "postgres",
		"database.uri":               "",
		"database.scheme":            "",
		"database.host":              "",
		"database.port":              5432,
		"database.schema_name":       "public",
		"database.ssl_mode":          "disable",
		"database.max_idle_conns":    10,
		"database.max_open_conns":    100,
		"database.conn_max_lifetime": "300s", // in seconds

		// Memory
		"memory.enabled":    true,
		"memory.limit":      0, // 0 = tidak dibatasi (dalam MB)
		"memory.expires_in": "180s",

		// Redis
		"redis.host": "",
		"redis.port": 6379,
		"redis.db":   0,

		// Kafka
		"kafka.enabled":      false,
		"kafka.brokers":      []string{},
		"kafka.group_id":     "",
		"kafka.topics":       []string{},
		"kafka.offset_reset": "earliest",

		// PubSub
		"pubsub.project_id":     "",
		"pubsub.topic":          "",
		"pubsub.subscription":   "",
		"pubsub.no_credentials": false,
		"pubsub.credentials":    "",

		// PubSub Consumer
		"pubsub.consumer.maxmessages":             10,
		"pubsub.consumer.sleeptime":               "5s",
		"pubsub.consumer.acktimeout":              "60s",
		"pubsub.consumer.retrycount":              3,
		"pubsub.consumer.retrydelay":              "1s",
		"pubsub.consumer.flowcontrol.enabled":     true,
		"pubsub.consumer.flowcontrol.maxmessages": 1000,
		"pubsub.consumer.flowcontrol.maxbytes":    1000000, // 1M

		// PubSub Producer
		"pubsub.producer.enableordering": false,
		"pubsub.producer.batchsize":      100,
		"pubsub.producer.attributes":     make(map[string]string),
	}
}
