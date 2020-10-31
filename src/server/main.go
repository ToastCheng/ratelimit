package main

import (
	"log"
	"os"
	"ratelimit"
	"strconv"
)

func getEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		vi, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("%s should be integer", key)
		}
		return vi
	}
	return fallback
}

func main() {
	c := ratelimit.Config{
		Port:   8080,
		Limit:  int(getEnv("RL_LIMIT", 60)),
		Window: int(getEnv("RL_WINDOW", 60)),
	}
	if err := c.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	s := ratelimit.New(c)
	log.Printf("start serving on %s..", s.Addr)
	log.Fatal(s.ListenAndServe())
}
