package ratelimit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigPort(t *testing.T) {
	c := Config{
		Port:   999999,
		Limit:  60,
		Window: 60,
	}
	err := c.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, PortValueError, err)
	}
}

func TestConfigLimit(t *testing.T) {
	c := Config{
		Port:   8080,
		Limit:  -1,
		Window: 60,
	}
	err := c.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, LimitValueError, err)
	}
}

func TestConfigWindow(t *testing.T) {
	c := Config{
		Port:   8080,
		Limit:  60,
		Window: -1,
	}
	err := c.Validate()
	if assert.Error(t, err) {
		assert.Equal(t, WindowValueError, err)
	}
}
