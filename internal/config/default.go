package config

var defaultConfig = Config{
	HttpServer: HttpServer{
		Port:                   8080,
		ShutdownTimeoutSeconds: 30,
	},
	Postgres: Postgres{
		Port:                  5432,
		DefaultTimeoutSeconds: 30,
	},
	StaleMailingEntryRemover: StaleMailingEntryRemover{
		StalenessThresholdSeconds: 2 * 60 * 60,
	},
}
