package config

// NatsConfig holds NATS connection configuration
type NatsConfig struct {
	URL string
}

// NewNatsConfig creates a new NATS configuration
func NewNatsConfig(url string) *NatsConfig {
	return &NatsConfig{
		URL: url,
	}
}

// GetURL returns the NATS URL
func (c *NatsConfig) GetURL() string {
	return c.URL
}
