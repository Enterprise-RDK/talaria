package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/xmidt-org/webpa-common/v2/device"
)

func TestDeviceMetadataMiddlewareKeepsTheAddressOfTheConnection(t *testing.T) {
	tests := []struct {
		description string
		remoteAddr  string
		expected    string
	}{
		{description: "IPv4 with a port", remoteAddr: "203.0.113.7:51234", expected: "203.0.113.7"},
		{description: "IPv6 with a port", remoteAddr: "[2001:db8::1]:443", expected: "2001:db8::1"},
		{description: "an address without a port", remoteAddr: "203.0.113.7", expected: "203.0.113.7"},
		{description: "no address", remoteAddr: "", expected: ""},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			var metadata *device.Metadata
			handler := DeviceMetadataMiddleware(func(context.Context) log.Logger { return nil })(
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
					var ok bool
					metadata, ok = device.GetDeviceMetadata(r.Context())
					assert.True(ok)
				}),
			)

			request := httptest.NewRequest(http.MethodGet, "/api/v2/device", nil)
			request.RemoteAddr = tc.remoteAddr
			handler.ServeHTTP(httptest.NewRecorder(), request)

			if assert.NotNil(metadata) {
				if tc.expected == "" {
					assert.Nil(metadata.Load(RemoteIPMetadataKey))
				} else {
					assert.Equal(tc.expected, metadata.Load(RemoteIPMetadataKey))
				}
			}
		})
	}
}
