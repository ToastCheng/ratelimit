package ratelimit

import "errors"

// Config defines the necessary configuration for ratelimit server.
type Config struct {
	Port   int // port to listen.
	Limit  int // request number allowed in [now-window, now]
	Window int // time interval for counting the incoming request
}

var PortValueError error = errors.New("port should be in the range [1024, 65353]")
var LimitValueError error = errors.New("limit should be greater than 0")
var WindowValueError error = errors.New("window should be greater than 0")

// Validate validates the value of Config.
func (c Config) Validate() error {
	if c.Port < 1024 || c.Port > 65353 {
		return PortValueError
	}
	if c.Limit < 0 {
		return LimitValueError
	}
	if c.Window < 0 {
		return WindowValueError
	}
	return nil
}
