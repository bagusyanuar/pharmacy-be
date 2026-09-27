package bootstrap_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/bagusyanuar/pharmacy-be/internal/bootstrap"
	"github.com/bagusyanuar/pharmacy-be/internal/config"
)

func TestDocsEndpoints(t *testing.T) {
	cfg := &config.Config{
		AppEnv: "test",
	}
	logger := zap.NewNop()

	app := bootstrap.SetupApp(cfg, nil, logger)

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{
			name:         "Health Check",
			url:          "/health",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Swagger UI HTML Page",
			url:          "/docs",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Swagger Redirect",
			url:          "/swagger",
			expectedCode: http.StatusFound, // 302 redirect to /docs
		},
		{
			name:         "API Docs Redirect",
			url:          "/api-docs",
			expectedCode: http.StatusFound, // 302 redirect to /docs
		},
		{
			name:         "Master YAML Spec",
			url:          "/docs/swagger.yaml",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Auth Module YAML Spec",
			url:          "/docs/modules/auth.yaml",
			expectedCode: http.StatusOK,
		},
		{
			name:         "User Module YAML Spec",
			url:          "/docs/modules/user.yaml",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("unexpected error making test request: %v", err)
			}
			if resp.StatusCode != tc.expectedCode {
				t.Errorf("url %s: expected status %d, got %d", tc.url, tc.expectedCode, resp.StatusCode)
			}
		})
	}
}
