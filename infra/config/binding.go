package config

func (c *Config) SetEnvBindings() map[string]string {
	return map[string]string{
		// App
		"app.name":                            "APP_NAME",
		"app.version":                         "APP_VERSION",
		"app.environment":                     "APP_ENVIRONMENT",
		"app.features.logging":                "APP_FEATURES_LOGGING",
		"app.features.recovery":               "APP_FEATURES_RECOVERY",
		"app.features.tracing":                "APP_FEATURES_TRACING",
		"app.features.metrics":                "APP_FEATURES_METRICS",
		"app.features.profiling":              "APP_FEATURES_PROFILING",
		"app.logging.level":                   "APP_LOGGING_LEVEL",
		"app.logging.format":                  "APP_LOGGING_FORMAT",
		"app.logging.output":                  "APP_LOGGING_OUTPUT",
		"app.logging.remote.driver":           "APP_LOGGING_REMOTE_DRIVER",
		"app.logging.remote.uri":              "APP_LOGGING_REMOTE_URI",
		"app.logging.remote.sample_rate":      "APP_LOGGING_REMOTE_SAMPLE_RATE",
		"app.logging.remote.structured":       "APP_LOGGING_REMOTE_STRUCTURED",
		"app.logging.remote.default_tags":     "APP_LOGGING_REMOTE_DEFAULT_TAGS",
		"app.logging.remote.default_contexts": "APP_LOGGING_REMOTE_DEFAULT_CONTEXTS",
		"app.security_headers":                "APP_SECURITY_HEADERS",
		"app.additional_headers":              "APP_ADDITIONAL_HEADERS",
		"app.cors.allow_origins":              "APP_CORS_ALLOW_ORIGINS",
		"app.cors.allow_methods":              "APP_CORS_ALLOW_METHODS",
		"app.cors.allow_headers":              "APP_CORS_ALLOW_HEADERS",
		"app.cors.expose_headers":             "APP_CORS_EXPOSE_HEADERS",
		"app.cors.allow_credentials":          "APP_CORS_ALLOW_CREDENTIALS",
		"app.cors.max_age":                    "APP_CORS_MAX_AGE",
		"app.rate_limit.enabled":              "APP_RATE_LIMIT_ENABLED",
		"app.rate_limit.max":                  "APP_RATE_LIMIT_MAX",
		"app.rate_limit.window":               "APP_RATE_LIMIT_WINDOW",
		"app.rate_limit.backend":              "APP_RATE_LIMIT_BACKEND",
		"app.rate_limit.fail_open":            "APP_RATE_LIMIT_FAIL_OPEN",
		"app.module.base_path":                "APP_MODULE_BASE_PATH",
		"app.module.disabled":                 "APP_MODULE_DISABLED",

		// Server
		"server.host":          "SERVER_HOST",
		"server.port":          "SERVER_PORT",
		"server.path":          "SERVER_PATH",
		"server.read_timeout":  "SERVER_READ_TIMEOUT",
		"server.write_timeout": "SERVER_WRITE_TIMEOUT",
		"server.proxy_header":  "SERVER_PROXY_HEADER",

		// Auth
		"auth.directory":             "AUTH_DIRECTORY",
		"auth.control":               "AUTH_CONTROL",
		"auth.store":                 "AUTH_STORE",
		"auth.type":                  "AUTH_TYPE",
		"auth.secret_key":            "AUTH_SECRET_KEY",
		"auth.session.backend":       "AUTH_SESSION_BACKEND",
		"auth.session.expires_in":    "AUTH_SESSION_EXPIRES_IN",
		"auth.session.refresh_in":    "AUTH_SESSION_REFRESH_IN",
		"auth.session.content_type":  "AUTH_SESSION_CONTENT_TYPE",
		"auth.session.public_client": "AUTH_SESSION_PUBLIC_CLIENT",
		"auth.session.cookie_name":   "AUTH_SESSION_COOKIE_NAME",
		"auth.session.username_key":  "AUTH_SESSION_USERNAME_KEY",
		"auth.session.password_key":  "AUTH_SESSION_PASSWORD_KEY",
		"auth.api_key_header":        "AUTH_API_KEY_HEADER",
		"auth.api_key_name":          "AUTH_API_KEY_NAME",

		// Database
		"database.driver":            "DATABASE_DRIVER",
		"database.uri":               "DATABASE_URI",
		"database.scheme":            "DATABASE_SCHEME",
		"database.host":              "DATABASE_HOST",
		"database.port":              "DATABASE_PORT",
		"database.user":              "DATABASE_USER",
		"database.password":          "DATABASE_PASSWORD",
		"database.name":              "DATABASE_NAME",
		"database.schema_name":       "DATABASE_SCHEMA_NAME",
		"database.ssl_mode":          "DATABASE_SSL_MODE",
		"database.max_open_conns":    "DATABASE_MAX_OPEN_CONNS",
		"database.max_idle_conns":    "DATABASE_MAX_IDLE_CONNS",
		"database.conn_max_lifetime": "DATABASE_CONN_MAX_LIFETIME",

		// Memory
		"memory.enabled":    "MEMORY_ENABLED",
		"memory.limit":      "MEMORY_LIMIT",
		"memory.expires_in": "MEMORY_EXPIRES_IN",

		// Redis
		"redis.host":     "REDIS_HOST",
		"redis.port":     "REDIS_PORT",
		"redis.password": "REDIS_PASSWORD",
		"redis.db":       "REDIS_DB",

		// Kafka
		"kafka.enabled":      "KAFKA_ENABLED",
		"kafka.brokers":      "KAFKA_BROKERS",
		"kafka.group_id":     "KAFKA_GROUP_ID",
		"kafka.topics":       "KAFKA_TOPICS",
		"kafka.offset_reset": "KAFKA_AUTO_OFFSET_RESET",

		// PubSub
		"pubsub.project_id":     "PUBSUB_PROJECT_ID",
		"pubsub.topic":          "PUBSUB_TOPIC",
		"pubsub.subscription":   "PUBSUB_SUBSCRIPTION",
		"pubsub.no_credentials": "PUBSUB_NO_CREDENTIALS",
		"pubsub.credentials":    "PUBSUB_CREDENTIALS",

		// PubSub Google Credentials
		"pubsub.credentials_data.type":                        "PUBSUB_CREDENTIALS_DATA_TYPE",
		"pubsub.credentials_data.project_id":                  "PUBSUB_CREDENTIALS_DATA_PROJECT_ID",
		"pubsub.credentials_data.private_key_id":              "PUBSUB_CREDENTIALS_DATA_PRIVATE_KEY_ID",
		"pubsub.credentials_data.private_key":                 "PUBSUB_CREDENTIALS_DATA_PRIVATE_KEY",
		"pubsub.credentials_data.client_email":                "PUBSUB_CREDENTIALS_DATA_CLIENT_EMAIL",
		"pubsub.credentials_data.client_id":                   "PUBSUB_CREDENTIALS_DATA_CLIENT_ID",
		"pubsub.credentials_data.auth_uri":                    "PUBSUB_CREDENTIALS_DATA_AUTH_URI",
		"pubsub.credentials_data.token_uri":                   "PUBSUB_CREDENTIALS_DATA_TOKEN_URI",
		"pubsub.credentials_data.auth_provider_x509_cert_url": "PUBSUB_CREDENTIALS_DATA_AUTH_PROVIDER_X509_URL",
		"pubsub.credentials_data.client_x509_cert_url":        "PUBSUB_CREDENTIALS_DATA_CLIENT_X509_URL",
		"pubsub.credentials_data.universe_domain":             "PUBSUB_CREDENTIALS_DATA_UNIVERSE_DOMAIN",

		// PubSub Consumer
		"pubsub.consumer.maxmessages":             "PUBSUB_CONSUMER_MAXMESSAGES",
		"pubsub.consumer.sleeptime":               "PUBSUB_CONSUMER_SLEEPTIME",
		"pubsub.consumer.acktimeout":              "PUBSUB_CONSUMER_ACKTIMEOUT",
		"pubsub.consumer.retrycount":              "PUBSUB_CONSUMER_RETRYCOUNT",
		"pubsub.consumer.retrydelay":              "PUBSUB_CONSUMER_RETRYDELAY",
		"pubsub.consumer.flowcontrol.enabled":     "PUBSUB_CONSUMER_FLOWCONTROL_ENABLED",
		"pubsub.consumer.flowcontrol.maxmessages": "PUBSUB_CONSUMER_FLOWCONTROL_MAXMESSAGES",
		"pubsub.consumer.flowcontrol.maxbytes":    "PUBSUB_CONSUMER_FLOWCONTROL_MAXBYTES",

		// PubSub Producer
		"pubsub.producer.enableordering": "PUBSUB_PRODUCER_ENABLEORDERING",
		"pubsub.producer.batchsize":      "PUBSUB_PRODUCER_BATCHSIZE",
		"pubsub.producer.attributes":     "PUBSUB_PRODUCER_ATTRIBUTES",
	}
}
