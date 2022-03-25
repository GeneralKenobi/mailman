package config

// Config aggregates configuration for modules.
type Config struct {
	HttpServer HttpServer `json:"httpServer"`
	Postgres   Postgres   `json:"postgres"`
}

// HttpServer stores configuration for the HTTP server.
type HttpServer struct {
	Port                   int `json:"port"`                   // Port to listen on
	ShutdownTimeoutSeconds int `json:"shutdownTimeoutSeconds"` // Graceful shutdown time
}

// Postgres stores configuration for connecting to the postgres database.
type Postgres struct {
	Host                  string `json:"host"`                  // DB server host, e.g. my-postgres.com or 10.101.146.170
	Port                  int    `json:"port"`                  // Port the DB is listening on
	User                  string `json:"user"`                  // User to use for connecting to the DB
	Password              string `json:"password"`              // Password authenticating User
	Database              string `json:"database"`              // Database to use
	DefaultTimeoutSeconds int    `json:"defaultTimeoutSeconds"` // Default timeout to use for DB operations if none is specified
}
