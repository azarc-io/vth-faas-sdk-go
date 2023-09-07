package http

import (
	"net/http"
	"time"
)

const defaultClientTimeout = 30 * time.Second

func GetDefaultClient() *http.Client {
	return &http.Client{Timeout: defaultClientTimeout}
}
