package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMetricsMiddlewareExportsHTTPMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	router := gin.New()
	router.Use(metrics.Middleware())
	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/failure", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})
	router.GET("/metrics", gin.WrapH(metrics.Handler()))

	for _, path := range []string{"/tasks/123", "/failure", "/missing"} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("GET /metrics status = %d, want %d", response.Code, http.StatusOK)
	}

	body := response.Body.String()
	for _, want := range []string{
		`devboard_http_requests_total{method="GET",route="/tasks/:id",status="200"} 1`,
		`devboard_http_request_duration_seconds_count{method="GET",route="/tasks/:id",status="200"} 1`,
		`devboard_http_requests_total{method="GET",route="/failure",status="500"} 1`,
		`devboard_http_requests_total{method="GET",route="unmatched",status="404"} 1`,
		"devboard_http_requests_in_flight 0",
		"go_goroutines",
		"process_cpu_seconds_total",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("metrics output missing %q", want)
		}
	}
	if strings.Contains(body, `route="/tasks/123"`) {
		t.Error("metrics output contains a per-task route label; expected route template")
	}
	if strings.Contains(body, `route="/metrics"`) {
		t.Error("metrics scrape endpoint should not be counted as an application request")
	}
}
