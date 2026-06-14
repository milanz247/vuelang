package models

import "time"

// Greeting is a simple response model — no DB table needed.
type Greeting struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// NewGreeting returns a fresh greeting with the current timestamp.
func NewGreeting() *Greeting {
	return &Greeting{
		Message:   "Welcome to Vuelang — Go + Vue 3, one binary. 🚀",
		Timestamp: time.Now().Format("2006-01-02 03:04:05 PM"),
	}
}
