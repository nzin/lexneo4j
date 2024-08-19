package config

// Config is the whole configuration of the app
var Config = struct {
	// Host - golang-skeleton server host
	Host string `env:"HOST" envDefault:"localhost"`
	// Port - golang-skeleton server port
	Port int `env:"PORT" envDefault:"18000"`

	// LogrusLevel sets the logrus logging level
	LogrusLevel string `env:"LEXNEO4J_LOGRUS_LEVEL" envDefault:"info"`
	// LogrusFormat sets the logrus logging formatter
	// Possible values: text, json
	LogrusFormat string `env:"LEXNEO4J_LOGRUS_FORMAT" envDefault:"json"`

	// MiddlewareVerboseLoggerEnabled - to enable the negroni-logrus logger for all the endpoints useful for debugging
	MiddlewareVerboseLoggerEnabled bool `env:"LEXNEO4J_MIDDLEWARE_VERBOSE_LOGGER_ENABLED" envDefault:"true"`
	// MiddlewareVerboseLoggerExcludeURLs - to exclude urls from the verbose logger via comma separated list
	MiddlewareVerboseLoggerExcludeURLs []string `env:"LEXNEO4J_MIDDLEWARE_VERBOSE_LOGGER_EXCLUDE_URLS" envDefault:"" envSeparator:","`
	// MiddlewareGzipEnabled - to enable gzip middleware
	MiddlewareGzipEnabled bool `env:"LEXNEO4J_MIDDLEWARE_GZIP_ENABLED" envDefault:"true"`

	// neo4j configuration
	Neo4jURL      string `env:"NEO4J_URL" envDefault:"bolt://neo4j:7687/neo4j"`
	Neo4jUsername string `env:"NEO4J_USERNAME" envDefault:"neo4j"`
	Neo4jPassword string `env:"NEO4J_PASSWORD" envDefault:"password"`
}{}
