package logger

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func TestCustomLogFormatter(t *testing.T) {
	// Create a sample gin.LogFormatterParams
	sampleParams := gin.LogFormatterParams{
		Request: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path: "/test",
			},
			Header: http.Header{
				"User-Agent": []string{"test_agent"},
			},
		},
		TimeStamp:    time.Now(),
		Latency:      time.Duration(100 * time.Millisecond),
		ClientIP:     "127.0.0.1",
		StatusCode:   200,
		ErrorMessage: "",
	}

	var color string
	switch {
	case sampleParams.StatusCode >= 500:
		color = red
	case sampleParams.StatusCode >= 400:
		color = yellow
	default:
		color = green
	}

	// Call the CustomLogFormatter function
	log := CustomLogFormatter(sampleParams)
	// Assert that the returned string has the expected format
	expectedFormat := "[LOGGER] [%s] \033[%sm %d \033[0m %s %s %d %s %s %s\n"

	assert.Equal(t, fmt.Sprintf(expectedFormat,
		sampleParams.TimeStamp.Format("2006-01-02 15:04:05"),
		color,
		sampleParams.StatusCode,
		sampleParams.Method,
		sampleParams.Path,
		sampleParams.Latency,
		sampleParams.ClientIP,
		sampleParams.ErrorMessage,
		sampleParams.Request.UserAgent(),
	), log)
}
