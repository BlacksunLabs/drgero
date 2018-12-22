package event

// Event holds information about an event
type Event struct {
	Message   string `json:"message"`
	Host      string `json:"host"`
	UserAgent string `json:"userAgent"`
}
