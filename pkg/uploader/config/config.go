// Package config manages the configuration options for the uploader.
package config

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDatabaseNotConfigured       = errors.New("database connection is not configured")
	ErrVisitorQueryNotConfigured   = errors.New("visitor query is not configured")
	ErrD2DConnectionFailed         = errors.New("unable to connect to https://integration.door2doc.net; make sure that the firewall allows access")
	ErrD2DCredentialsNotConfigured = errors.New("username/password is not configured")
)

// ValidationResult contains the results of validating the current configuration.
type ValidationResult struct {
	DatabaseConnection error
	VisitorQuery       error
	D2DConnection      error
	D2DCredentials     error
}

// IsValid returns true if all possible validation errors are nil.
func (v *ValidationResult) IsValid() bool {
	// todo
	return true
}

// Configuration contains the configuration options for the service.
type Configuration struct {
	// username to connect to the d2d upload service
	Username string
	// password to connect to the d2d upload service
	Password string
	// DSN for the database connection to retrieve visitor information from
	DSN string
	// Query to execute to retrieve visitor information
	Query string
	// Set to true if the service should be active
	Active bool
	// Pause between runs
	Interval time.Duration
}

// Load loads the configuration from a well-known location. It does not give an error if the configuration
// does not exist.
func Load(ctx context.Context) (*Configuration, error) {
	// todo
	return &Configuration{
		Active:   true,
		Interval: time.Minute,
	}, nil
}

// Validate validates the configuration and returns the results of those checks.
func (c *Configuration) Validate(ctx context.Context) *ValidationResult {
	// todo
	return &ValidationResult{}
}

// Save stores the latest configuration values to a well-known location.
func (c *Configuration) Save() error {
	// todo
	return nil
}

// Refresh ensures that the configuration is the latest version saved.
func (c *Configuration) Refresh() error {
	// todo
	return nil
}
