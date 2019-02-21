package server

import (
	"errors"
	"time"
)

var (
	// errNilServerConfig is returned on nil ServerConfig
	errNilServerConfig = errors.New("nil ServerConfig")
)

// ServerConfig hold server configuration.
type ServerConfig struct {
	EtcdURLs       []string
	LockTimeout    time.Duration
	SemaphoreSlots uint64
}
