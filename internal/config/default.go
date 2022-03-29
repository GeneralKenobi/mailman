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
		StalenessThresholdSeconds: 5 * 60, // 5 minutes
	},
	MailingEntryCleanupJob: MailingEntryCleanupJob{
		PeriodSeconds: 60 * 60, // 1 hour
	},
}
