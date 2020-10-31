package ratelimit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryOne(t *testing.T) {
	s := New(Config{
		Port:   8080,
		Limit:  5,
		Window: 5,
	})
	ts := httptest.NewServer(s.Handler)

	c := http.Client{}
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	resp, err := c.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.NoError(t, err)
	assert.Equal(t, "Current request: 1", bodyStr)
}

func TestQueryRepeat(t *testing.T) {
	s := New(Config{
		Port:   8080,
		Limit:  1,
		Window: 1,
	})
	ts := httptest.NewServer(s.Handler)

	c := http.Client{}
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	for i := 0; i < 4; i++ {
		resp, err := c.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		assert.NoError(t, err)
		assert.Equal(t, "Current request: 1", bodyStr)
		time.Sleep(1 * time.Second)
	}
}

func TestQueryFromDifferentIP(t *testing.T) {
	s := New(Config{
		Port:   8080,
		Limit:  1,
		Window: 1,
	})
	ts := httptest.NewServer(s.Handler)

	c := http.Client{}
	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		assert.NoError(t, err)
		req.Header.Add("X-Forwarded-For", fmt.Sprintf("127.0.0.%d", i))
		resp, err := c.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		assert.NoError(t, err)
		assert.Equal(t, "Current request: 1", bodyStr)
	}
}
