package ping_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sfshf/gonoweb/internal/app/web/ginsrv"
	"github.com/stretchr/testify/assert"
)

// go test -v -count=1 -timeout 30s -run ^TestPing$ github.com/sfshf/gonoweb/internal/app/web/ginsrv/ping
func TestPing(t *testing.T) {
	r, err := ginsrv.InitGin(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "PONG")
}
