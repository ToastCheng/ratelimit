package ratelimit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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

func TestQueryConcurrentRepeat(t *testing.T) {
	s := New(Config{
		Port:   8080,
		Limit:  3,
		Window: 1,
	})
	ts := httptest.NewServer(s.Handler)

	c := http.Client{}
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	var r200 int32 = 0
	var r429 int32 = 0
	makeReq := func(wg *sync.WaitGroup, r200, r429 *int32) {
		resp, err := c.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			atomic.AddInt32(r200, 1)
		} else if resp.StatusCode == 429 {
			atomic.AddInt32(r429, 1)
		}

		wg.Done()
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go makeReq(wg, &r200, &r429)
	}

	wg.Wait()
	assert.Equal(t, int32(3), r200)
	assert.Equal(t, int32(2), r429)
}
