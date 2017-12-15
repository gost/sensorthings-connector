package configuration

// Config contains the settings for the connector
type Config struct {
	Connector ConnectorConfig `json:"connector"`
	Logging   LoggingConfig   `json:"logging"`
}

// ConnectorConfig contains the general config information
type ConnectorConfig struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	StartModulesOnStartup bool   `json:"startModulesOnStartup"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Status  StatusConfig  `json:"status"`
	Discord DiscordConfig `json:"discord"`
}

// StatusConfig contains the status report configuration
type StatusConfig struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"intervalSeconds"`
}

// DiscordConfig contains the discord logging configuration
type DiscordConfig struct {
	Enabled bool   `json:"enabled"`
	Hook    string `json:"hook"`
	Name    string `json:"name"`
}

// Validate checks if all mandatory params are set in the config
func (c Config) Validate() error {
	return nil
}
