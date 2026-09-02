package qzone_api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/infrastructure/config"
	"strings"
	"testing"
)

func TestGetLoginQRCodeValidatesResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		setCookie  bool
		wantError  string
	}{
		{name: "success", statusCode: http.StatusOK, setCookie: true},
		{name: "http error", statusCode: http.StatusBadGateway, setCookie: true, wantError: "502 Bad Gateway"},
		{name: "missing qrsig", statusCode: http.StatusOK, wantError: "缺少 qrsig"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.setCookie {
					http.SetCookie(w, &http.Cookie{Name: "qrsig", Value: "test-signature"})
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("png-data"))
			}))
			defer server.Close()

			cfg := &config.Config{}
			cfg.QzoneAPI.QRCodeURL = server.URL
			client := &qzoneAPIClient{httpClient: server.Client(), config: cfg}
			body, signature, err := client.GetLoginQRCode()

			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("GetLoginQRCode() error = %v, want substring %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetLoginQRCode() unexpected error: %v", err)
			}
			if string(body) != "png-data" || signature != "test-signature" {
				t.Fatalf("GetLoginQRCode() = (%q, %q), want (%q, %q)", body, signature, "png-data", "test-signature")
			}
		})
	}
}

func TestCheckLoginStatusUsesProtocolCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		response   string
		wantStatus entity.LoginStatus
		wantError  bool
	}{
		{name: "waiting", response: `ptuiCB('66','0','','0','localized text','')`, wantStatus: entity.LoginStatusWaiting},
		{name: "scanning", response: `ptuiCB('67','0','','0','localized text','')`, wantStatus: entity.LoginStatusScanning},
		{name: "expired", response: `ptuiCB('65','0','','0','localized text','')`, wantStatus: entity.LoginStatusExpired},
		{name: "success", response: `ptuiCB('0','0','https://example.test','0','localized text','')`, wantStatus: entity.LoginStatusSuccess},
		{name: "unknown", response: `ptuiCB('999','0','','0','localized text','')`, wantStatus: entity.LoginStatusWaiting, wantError: true},
		{name: "malformed", response: `not a ptui response`, wantStatus: entity.LoginStatusWaiting, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{}
			cfg.QzoneAPI.LoginURL = server.URL
			client := &qzoneAPIClient{httpClient: server.Client(), config: cfg}
			status, response, err := client.CheckLoginStatus("signature")

			if (err != nil) != tt.wantError {
				t.Fatalf("CheckLoginStatus() error = %v, wantError %t", err, tt.wantError)
			}
			if status != tt.wantStatus {
				t.Fatalf("CheckLoginStatus() status = %v, want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == entity.LoginStatusSuccess && response != tt.response {
				t.Fatalf("CheckLoginStatus() response = %q, want %q", response, tt.response)
			}
		})
	}
}

func TestReadResponseBodyRejectsOversizedPayload(t *testing.T) {
	t.Parallel()

	resp := &http.Response{Body: io.NopCloser(strings.NewReader(strings.Repeat("x", maxResponseBodySize+1)))}
	if _, err := readResponseBody(resp); err == nil {
		t.Fatal("readResponseBody() expected oversized response error")
	}
}
