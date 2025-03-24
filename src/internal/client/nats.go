package client

import (
	"github.com/mjmhtjain/knime/src/config"
	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn = nil

// NewNatsClient returns a singleton NATS connection
func NewNatsClient(config *config.NatsConfig) (*nats.Conn, error) {
	if natsConn != nil {
		return natsConn, nil
	}

	var err error
	natsConn, err = nats.Connect(config.GetURL())
	if err != nil {
		return nil, err
	}

	// Test the connection
	if !natsConn.IsConnected() {
		return nil, nats.ErrConnectionClosed
	}

	return natsConn, nil
}

// CloseNatsConnection closes the NATS connection
func CloseNatsConnection() {
	if natsConn != nil {
		natsConn.Close()
		natsConn = nil
	}
}
